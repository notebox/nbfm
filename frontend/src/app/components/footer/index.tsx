import type {Navigator, ModeState, FileInfo} from "@/domain"

import React from "react"
import {NodeMenu} from "./menu"
import "./styles.scss"
import {IsEscape} from "@/app/utils"

const comp = ({nav, mode}: {nav: Navigator; mode: ModeState; current: string}) => {
  if (mode === "searching") {
    return (
      <footer>
        <Search nav={nav} />
      </footer>
    )
  }

  const idx = nav.idx

  return (
    <footer>
      {idx < 0 ? null : <Status file={nav.dir.files[idx]} index={idx} total={nav.dir.files.length} />}
      {mode === "node-menu" ? <NodeMenu nav={nav} /> : null}
    </footer>
  )
}

const Search = ({nav}: {nav: Navigator}) => {
  return (
    <div className="ui-search">
      <div>/</div>
      <input autoFocus={true} onKeyDown={e => {
        if (e.key === "Enter") {
          e.preventDefault()
          e.stopPropagation()
          nav.search((e.target as any).value)
        } else if (IsEscape(e)) {
          e.preventDefault()
          e.stopPropagation()
          nav.setMode(null)
        }
      }} />
    </div>
  )
}

const Status = ({file, index, total}: {file: FileInfo; index: number; total: number}) => {
  return (
    <div className="ui-status">
      <div className="ui-mode">{file.mode}</div>
      <div>{file.username}</div>
      <div>{file.groupName}</div>
      <div>{file.size}</div>
      <div>{file.modTime}</div>
      <div className="ui-idx">{index + 1}/{total}</div>
    </div>
  )
}

export const Footer = React.memo(comp, (prev, next) => {
  return prev.current === next.current && prev.mode === next.mode
})