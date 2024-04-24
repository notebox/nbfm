import type {Navigator, PreviewInfo} from "@/domain"

import React from "react"
import {useRecoilValue} from "recoil"

import {Manual} from "../manual"
import {Dir} from "../dir"
import "./styles.scss"
import {NoteboxEditor} from "./editor"

const comp = ({nav}: {nav: Navigator}) => {
  const preview = useRecoilValue(nav.state.preview.atom)

  if (!preview.path) {
    return null
  }

  return <KeyedComp key={preview.path} nav={nav} preview={preview} />
}

const KeyedComp = ({nav, preview}: {nav: Navigator; preview: PreviewInfo}) => {
  if (preview.dir) {
    if (preview.path.endsWith(".note")) {
      return <div className="ui-preview-embed"><NoteboxEditor path={preview.path} /></div>
    }

    return <Dir dir={preview.dir} hidden={nav.hidden} />
  }

  switch (preview.type) {
  case "embed":
    return <div className="ui-preview-embed"><embed src={preview.path+"?"+Date.now()} type="application/pdf" /></div>
  case "image":
    return <div className="ui-preview-image"><img src={preview.path} alt={preview.path} /></div>
  case "video":
    return <div className="ui-preview-video"><video controls src={preview.path} autoPlay /></div>
  case "audio":
    return <div className="ui-preview-audio"><audio controls src={preview.path} autoPlay /></div>
  case "text":
    if (!preview.utf8) {
      return <Manual style={{opacity: "0.5"}} />
    }
    return <div className="ui-preview-text"><pre>{preview.utf8}</pre></div>
  case "error":
    return <div className="ui-preview-error">{preview.err}</div>
  default:
    return <div className="ui-preview-error">binary</div>
  }
}

export const Preview = React.memo(comp)