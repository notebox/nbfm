import type {Navigator, ModeState} from "@/domain"

import React from "react"

import {NBEditorToolbar} from "../preview/editor/toolbar"
import "./styles.scss"

const comp = ({nav, mode, homeDir, current}: {nav: Navigator, mode: ModeState; homeDir: string; current: string}) => {
  const path = current.startsWith(homeDir) ? current.replace(homeDir, "~") : current
  const frags = path.split(nav.separator)
  const focus = frags.pop()

  return (
    <header>
      <div className="ui-header-path"><span>{frags.join(nav.separator)+nav.separator}</span>{focus}</div>
      {!mode && current.endsWith(".note") ? <NBEditorToolbar nav={nav} /> : null}
    </header>
  )
}

export const Header = React.memo(comp, (prev, next) => {
  return prev.homeDir === next.homeDir && prev.current === next.current && prev.mode === next.mode
})