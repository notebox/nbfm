import type {DirInfo, FileInfo} from "@/domain"

import React, {useEffect, useRef} from "react"
import {nav} from "@/domain"
import "./styles.scss"

const comp = ({dir}: {dir: DirInfo; hidden: boolean}) => {
  const scrollable = useRef<HTMLDivElement>(null)
  useEffect(() => {
    const dom = scrollable.current
    if (!dom) return
    const selected = dom.querySelector("[data-selected=\"true\"]")
    if (!selected) return
    selected.scrollIntoView({block: "nearest"})
  }, [dir.selected, scrollable.current])

  return (
    <div className="ui-dir" ref={scrollable}>
      {dir.files.map((file: FileInfo) => (
        <div
          onClick={e => {
            e.preventDefault()
            e.stopPropagation()
            nav.select(dir.path, file.name)
          }}
          className="ui-item"
          data-editable={file.name.endsWith(".note")}
          data-dir={file.isDir}
          data-hidden={file.name.startsWith(".")}
          data-selected={dir.selected === file.name}
          key={file.name}>
          {file.name}
        </div>
      ))}
    </div>
  )
}

export const Dir = React.memo(comp, (prev, next) => {
  return prev.dir === next.dir
})