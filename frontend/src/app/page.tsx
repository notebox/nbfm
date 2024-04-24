
import {useEffect} from "react"
import {RecoilRoot} from "recoil"
import {useRecoilValue} from "recoil"

import {RecoilExternalStatePortal} from "@/domain/usecase/state/common"
import {nav} from "@/domain"
import {Dir} from "./components/dir"
import {ManualLayer} from "./components/manual"
import {Preview} from "./components/preview"
import {Header} from "./components/header"
import {Footer} from "./components/footer"
import "./styles.scss"

const App = () => (
  <RecoilRoot>
    <Content />
    <RecoilExternalStatePortal />
  </RecoilRoot>
)

const Content = () => {
  useRecoilValue(nav.state.redraw.atom)

  useEffect(() => {
    window.addEventListener("keydown", nav.onKeyDown)

    return () => {
      window.removeEventListener("keydown", nav.onKeyDown)
    }
  }, [window])

  if (!nav.initialized) {
    return <div style={{
      height: "100vh",
      width: "100vw",
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
      fontSize: "2em",
    }}>Loading...</div>
  }

  return (
    <div id="app">
      <Header nav={nav} mode={nav.mode} homeDir={nav.homeDir} current={nav.current} />
      <div id="layout-left">
        <Dir dir={nav.parent} hidden={nav.hidden} />
      </div>
      <div id="layout-center">
        <Dir dir={nav.dir} hidden={nav.hidden} />
      </div>
      <div id="layout-right">
        <Preview nav={nav} />
      </div>
      <Footer nav={nav} mode={nav.mode} current={nav.current} />
      {nav.mode === "manual" ? <ManualLayer nav={nav} /> : null}
    </div>
  )
}

export default App
