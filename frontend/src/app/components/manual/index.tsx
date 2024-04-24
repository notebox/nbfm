import type {Navigator} from "@/domain"

import {useEffect} from "react"
import "./styles.scss"

export const ManualLayer = ({nav}: {nav: Navigator}) => {
  useEffect(() => {
    const keymap = () => nav.setMode(null)
    window.addEventListener("keydown", keymap)
    return () => window.removeEventListener("keydown", keymap)
  }, [])

  return (
    <div id="ui-manual-layer">
      <div id="ui-dismiss">press any key to dismiss</div>
      <Manual />
    </div>
  )
}

export const Manual = ({style}: {style?: React.CSSProperties}) => (
  <div id="ui-manual" style={style}>
    <h1><b>NBFM</b> manual</h1>
    <h2>shortcuts</h2>
    <div className="ui-body">
      <dt>goto</dt><dd>(<span>H</span>)ome or (<span>~</span>) directory</dd>
      <dt>move</dt><dd><span>hjkl</span> or <span>←↓↑→</span></dd>
      <dt><span>gg</span></dt><dd>move to the top</dd>
      <dt><span>G</span></dt><dd>move to the bottom</dd>

      <dt><span>?</span></dt><dd>manual</dd>
      <dt><span>esc</span></dt><dd>dismiss current mode</dd>
      <dt><span>/</span></dt><dd>search</dd>
      <dt><span>.</span></dt><dd>toggle showing hidden dot files for mac</dd>
      <dt><span>m</span></dt>
      <dd>
        <div>
          menu
        </div>
        <div className="ui-content">
          <div>(<span>n</span>)ew note</div>
          <div>(<span>a</span>)add a child node</div>
          <div>(<span>m</span>)ove the current node</div>
          <div>(<span>d</span>)elete the current node</div>
          <div>(<span>c</span>)opy the current node</div>
          <div>copy (<span>p</span>)ath to clipboard</div>
          <div>(<span>l</span>)ist contents (only if .note dir)</div>
        </div>
      </dd>
      <div className="ui-not-yet">
        <dt><span>$</span></dt><dd>run system command in the current directory (<span>not yet support</span>)</dd>
      </div>
    </div>
  </div>
)