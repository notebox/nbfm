import type {Navigator} from "@/domain"

import {useEffect} from "react"
import {useRecoilValue} from "recoil"
import {DeleteFile} from "gojs/main/App"

export const Delete = ({nav, height}: {nav: Navigator, height: number}) => {
  return (
    <div className="ui-node-menu" style={{height}}>
      <header>Delete the current node</header>
      <Content nav={nav} />
    </div>
  )
}

const Content = ({nav}: {nav: Navigator}) => {
  const preview = useRecoilValue(nav.state.preview.atom)
  const notEmptyDir = !!preview.dir?.files.length

  if (notEmptyDir) {
    return (
      <div className="ui-content ui-for-delete">
        <div>STOP! Directory is not empty! To delete, type 'yes'</div>
        <div className="ui-confirm">
          {nav.current}:
          <input autoFocus={true} onKeyDown={e => {
            if (e.key === "Enter") {
              e.preventDefault()
              e.stopPropagation()
              if ((e.target as any).value !== "yes") {
                nav.setMode(null)
                return
              }
              DeleteFile(nav.current)
              nav.reload()
            }
          }} />
        </div>
      </div>
    )
  }

  useEffect(() => {
    const keymap = (event: KeyboardEvent) => {
      event.preventDefault()
      event.stopPropagation()

      switch (event.code) {
      case "KeyY":
        DeleteFile(nav.current)
        nav.reload()
        return
      default:
        nav.setMode(null)
        return
      }
    }
    window.addEventListener("keydown", keymap)
    return () => window.removeEventListener("keydown", keymap)
  }, [])

  return (
    <div className="ui-content">{nav.current} (yN):</div>
  )
}
