import type {RecoilState} from "recoil"

import {atom as genAtom} from "recoil"
import {setRecoilExternalState} from "./common"

export class EditorToolbarStateHandler {
  readonly atom: RecoilState<EditorToolbarState>
  private mode: EditorToolbarState = null

  constructor() {
    this.atom = genAtom<EditorToolbarState>({
      key: "editor-toolbar",
      default: null,
    })
  }

  set = (newState: EditorToolbarState) => {
    this.mode = newState
    setRecoilExternalState(this.atom, this.mode)
  }

  get = () => {
    return this.mode
  }
}

export type EditorToolbarState = "block" | "color" | null;