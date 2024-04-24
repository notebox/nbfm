import {WD, ReadDirFiles, Preview} from "gojs/main/App"

import {IsEscape} from "@/app/utils"
import {LRUSet} from "./lru"
import state from "../state"

export class Navigator {
  readonly state = state 
  separator: string
  homeDir: string
  hidden = false

  private lru: LRUSet<string>
  private _parent: DirInfo
  private _dir: DirInfo

  private defaultPath: string
  private _mode: ModeState

  constructor() {
    this.lru = new LRUSet(1024)
    this.separator = ""
    this.homeDir = ""
    this.defaultPath = ""
    this._parent = {path: "", files: [], selected: ""}
    this._dir = {path: "", files: [], selected: ""}
    this._mode = null
    this.initialize();
    (window as any).nbfm = this
  }

  get initialized() {
    return !!this.homeDir
  }

  get parent() {
    return this._parent
  }

  get dir() {
    return this._dir
  }

  get current() {
    return this._dir.path === this.separator ? this._dir.path + this._dir.selected : `${this._dir.path}${this.separator}${this._dir.selected}`
  }

  get idx() {
    return this._idx
  }

  get mode() {
    return this._mode
  }

  onKeyDown = (e: KeyboardEvent) => {
    if (this.mode) return
    if (!this.dir.path) return
    if (e.target && (e.target as any).tagName !== "BODY") return

    e.preventDefault()
    e.stopPropagation()
    if (IsEscape(e)) {
      this.setMode(null)
      return
    }

    switch (e.key) {
    case "?":
      if (!this.mode) {
        this.setMode("manual")
        return
      }
      break
    case "~": {
      this.select(this.homeDir)
      return
    }
    }
    switch (e.code) {
    case "Enter":
      this.focusNote()
      break
    case "Slash":
      this.setMode("searching")
      break
    case "Period":
      this.toggleHidden()
      break
    case "ArrowDown":
    case "KeyJ": {
      const idx = this.idx
      if (idx === this.dir.files.length - 1) return
      this.selectIDX(idx+1)
      return
    }
    case "ArrowUp":
    case "KeyK": {
      const idx = this.idx
      if (idx === 0) return
      this.selectIDX(idx-1)
      return
    }
    case "ArrowLeft":
    case "KeyH": {
      if (e.shiftKey) {
        this.select(this.defaultPath)
        return
      }
      if (!this._parent.path) return
      this.select(this._dir.path)
      return
    }
    case "ArrowRight":
    case "KeyL": {
      if (this.hasContents()) return
      this.listContents()
      return
    }
    case "KeyG": {
      if (e.shiftKey) {
        this.selectIDX(this.dir.files.length - 1)
        return
      }
      if (this.lastG && Date.now() - this.lastG < 500) {
        this.selectIDX(0)
        return
      } else {
        this.lastG = Date.now()
      }
      return
    }
    case "KeyN":
      e.shiftKey ? this.searchPrev() : this.searchNext()
      return
    case "KeyM":
      this.setMode("node-menu")
      return
    }
  }

  hasContents() {
    return this.current.endsWith(".note")
  }

  listContents() {
    const selected = this.selected(this.current)
    if (!selected) return
    this.select(`${this.current}${this.separator}${selected}`)
  }

  setMode(mode: ModeState) {
    if (this.mode === mode) return
    this._mode = mode
    this.state.redraw.set()
  }

  toggleHidden() {
    this.hidden = !this.hidden
    this.select(this.current)
  }

  search(query: string) {
    this.setMode(null)
    this.lastQ = query
    this.searchNext()
  }

  searchNext() {
    if (!this.lastQ) return
    const regex = new RegExp(this.lastQ, "i")
    for (let i = this.idx + 1; i < this.dir.files.length; i++) {
      const file = this.dir.files[i]
      if (!file) break
      if (regex.test(file.name)) {
        this.selectIDX(i)
        return
      }
    }
  }
  searchPrev() {
    if (!this.lastQ) return
    const regex = new RegExp(this.lastQ, "i")
    for (let i = this.idx - 1; i >= 0; i--) {
      const file = this.dir.files[i]
      if (!file) break
      if (regex.test(file.name)) {
        this.selectIDX(i)
        return
      }
    }
  }

  reload() {
    this.setMode(null)
    this.select(this.current)
  }

  async select(...paths: string[]) {
    let path = paths.join(this.separator)
    if (path.endsWith(this.separator)) {
      path = path.slice(0, -1)
    }
    const [dirPath, dirSelected] = splitFilename(path, this.separator)
    this.lru.add(dirPath, dirSelected)
    if (dirPath === this.separator) {
      this._select(await this.readDirFiles(dirPath))
      return
    } else {
      const [parentPath, parentSelected] = splitFilename(dirPath, this.separator)
      this.lru.add(parentPath, parentSelected)
      const [dir, parent] = await Promise.all([
        this.readDirFiles(dirPath),
        this.readDirFiles(parentPath),
      ])
      this._select(dir, parent)
    }
  }

  private selected(path: string) {
    return this.lru.get(path)
  }

  private focusNote() {
    if (!this.current.endsWith(".note")) return
    (window as any).nbExecutor(null, (nb: any) => {
      nb.editor.emitter.handler.focus()
    })
  }

  private async readDirFiles(path: string): Promise<DirInfo> {
    const files = await ReadDirFiles(path, this.hidden)
    let selected = this.selected(path)
    if (!selected && files[0]) {
      selected = files[0].name
      this.lru.add(path, selected)
    }
    return {path, files, selected}
  }

  private _select(dir: DirInfo, parent?: DirInfo) {
    const idx = dir.files.findIndex(f => f.name === dir.selected)
    if (idx === -1) {
      const selected = dir.files[this.idx] || dir.files[this.idx - 1] || dir.files[dir.files.length - 1]
      if (selected) {
        this.lru.add(dir.path, selected.name)
        this._select({...dir, selected: selected.name}, parent)
        return
      }
      dir.selected = ""
    }
    this._parent = parent ?? {path: "", files: [], selected: ""}
    this._dir = dir
    this._idx = idx
    this.setMode(null)
    this.setPreview(this.current)
    this.state.redraw.set()
  }

  private _idx = -1
  private lastG?: number
  private lastQ?: string

  private selectIDX(idx: number) {
    const selected = this.dir.files[idx].name
    this.lru.add(this.dir.path, selected)
    this._select({...this.dir, selected}, this.parent)
  }

  private async initialize() {
    const wd = await WD()
    const [home, cur] = await Promise.all([
      this.readDirFiles(wd.HomeDir),
      this.readDirFiles(wd.Path),
    ])
    this.separator = wd.Separator
    this.homeDir = firstPathOf(home, this.separator)
    this.defaultPath = firstPathOf(cur, this.separator)
    await this.select(this.defaultPath)
    this.setMode("manual")
  }

  private async setPreview(path: string) {
    if (path !== this.current) return
    const preview = await this.previewInfo(this.current)
    if (path !== this.current) return
    this.state.preview.set(preview)
  }

  private async previewInfo(path: string): Promise<PreviewInfo> {
    try {
      const preview = await Preview(path, this.hidden)
      if (preview.type === "dir") {
        let selected = this.selected(path)
        if (!selected && preview.dirFiles?.[0]) {
          selected = preview.dirFiles[0].name
          this.lru.add(path, selected)
        }
        return {
          path,
          dir: {
            path,
            files: preview.dirFiles ?? [],
            selected,
          },
          type: preview.type,
        }
      }
      return {...preview, type: preview.type as PreviewInfo["type"], path}
    } catch (e) {
      return {path, type: "error", err: String(e)}
    }
  }
}

const splitFilename = (path: string, separator: string): [string, string | undefined] => {
  const parts = path.split(separator)
  const filename = parts.pop()
  return [parts.join(separator) || separator, filename]
}

const firstPathOf = (dir: DirInfo, separator: string): string => {
  return dir.files.length ? `${dir.path}${separator}${dir.files[0].name}` : dir.path
}

export const nav = new Navigator()

export type PreviewInfo = {
  path: string;
  dir?: DirInfo;
  utf8?: string;
  type: "dir" | "text" | "image" | "video" | "audio" | "embed" | "unknown" | "error";
  err?: string;
};
export type FileInfo = {
  name:  string;
  rawSize:  number;
  isDir: boolean;

  mode: string;
  username: string;
  groupName: string;
  size: string;
  modTime: string;
}
export type DirInfo = {
  files: FileInfo[];
  path: string;
  selected: string;
}
export type ModeState = "manual" | "searching" | "node-menu" | null;