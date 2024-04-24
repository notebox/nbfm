import type {RecoilState} from "recoil"
import type {PreviewInfo} from "@/domain/usecase/nav"

import {atom as genAtom} from "recoil"
import {setRecoilExternalState} from "./common"

export class PreviewStateHandler {
  readonly atom: RecoilState<PreviewInfo>

  constructor() {
    this.atom = genAtom<PreviewInfo>({
      key: "preview",
      default: {
        path: "",
        type: "unknown",
      },
    })
  }

  set(newState: PreviewInfo) {
    setRecoilExternalState(this.atom, newState)
  }
}
