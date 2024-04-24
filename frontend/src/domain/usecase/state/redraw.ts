import type {RecoilState} from "recoil"

import {atom as genAtom} from "recoil"
import {setRecoilExternalState} from "./common"

export class RedrawStateHandler {
  readonly atom: RecoilState<number>

  constructor() {
    this.atom = genAtom<number>({
      key: "redraw",
      default: 0,
    })
  }

  set() {
    setRecoilExternalState(this.atom, cur => cur > 999_999_999 ? 0 : cur + 1)
  }
}
