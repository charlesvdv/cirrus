import type { Component } from 'solid-js'
import { lazy } from 'solid-js'
import { Routes, Route } from "@solidjs/router"

import Layout from './components/Layout'

const App: Component = () => {
  return (
    <Layout>
      <Routes>
      </Routes>
    </Layout>
  )
}

export default App
