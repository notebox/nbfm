import {NBError, NBConnected, NBContribute} from "gojs/main/App"
import {BrowserOpenURL} from "runtime/runtime"
import state from "@/domain/usecase/state"
import {IsEscape} from "@/app/utils"

class NBEditorMessageHandler {
  nbError(identifier: string, err: string): void {
    NBError(identifier, err)
  }

  nbConnected(identifier: string, connected: boolean): void {
    (window as any).nbExecutor(identifier, (nb: any) => nb.setTheme("black"))
    NBConnected(identifier, connected)
  }

  nbContribute(identifier: string, ctrbs: string): void {
    NBContribute(identifier, ctrbs)
  }

  nbInitiated(_identifier: string, _initiated: boolean): void {
  }

  nbSelected(_identifier: string, _selection: string): void {
    if (state.toolbar.get()) {
      state.toolbar.set(null)
    }
  }

  nbNavigate(_identifier: string, url: string): void {
    BrowserOpenURL(url)
  }

  nbUploadFile(_identifier: string, _json: string): void {
  }

  nbHaptic(_identifier: string, _haptic: boolean): void {
  }

  nbDraggingStartByTouch(_identifier: string, _dragging: boolean): void {
  }

  nbDraggingEndByTouch(_identifier: string, _bool: boolean): void {
  }

  nbSearchMode(_identifier: string, _bool: boolean): void {
  }

  nbOpenFile(_identifier: string, _json: string): void {
  }

  fileURL(_identifier: string, _src: string | undefined, _fileID: string | undefined): string {
    return ""
  }

  previewLink(_identifier: string, _url: string) {
    return
  }

  private selection: any
  keymap = (ctx: any, event: KeyboardEvent): boolean => {
    if (IsEscape(event)) {
      this.selection = ctx.editor.selection
      window.focus()
    }
    return false
  }

  focus(): void {
    const selection = this.selection
    delete this.selection;
    (window as any).nbExecutor(null, (nb: any) => {
      nb.focus()
      if (selection) nb.editor.select(selection)
    })
  }
}

export const nbEditorMessageHandler = new NBEditorMessageHandler()