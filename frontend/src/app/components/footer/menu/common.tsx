import type {Navigator} from "@/domain"

export const NodeExecTemplate = ({
  nav,
  height,
  header,
  description,
  defaultValue,
  suffix,
  exec,
}: {
  nav: Navigator,
  height: number,
  header: string,
  description: string,
  defaultValue: string,
  suffix?: string,
  exec: (path: string) => Promise<void>,
}) => {
  return (
    <div className="ui-node-menu" style={{height}}>
      <header>{header}</header>
      <div className="ui-content">
        <div>{description}</div>
        <div className="ui-confirm">
          <input autoFocus={true} onKeyDown={e => {
            if (e.key === "Enter") {
              e.preventDefault()
              e.stopPropagation()
              const path = (e.target as any).value
              if (!path || path === defaultValue) {
                nav.setMode(null)
                return
              }
              const dst = suffix ? path + suffix : path
              exec(dst)
                .then(() => nav.select(dst))
                .catch(() => nav.setMode(null))
            }
          }} defaultValue={defaultValue} />
          {suffix ? <span>{suffix}</span> : null}
        </div>
      </div>
    </div>
  )
}
