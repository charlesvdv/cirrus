import { Component, createResource, Show, Suspense } from 'solid-js'
import { lazy } from 'solid-js'
import { Routes, Route, Navigate } from '@solidjs/router'

import { Client } from './lib/client'

import Layout from './components/Layout'
const GettingStarted = lazy(() => import('./pages/getting-started'))

const App: Component = () => {
  const client = new Client()

  const [instance] = createResource(() => client.getInstance())

  return (
    <Layout>
      <Suspense>
        {/* When a instance is not initialized, force the initialization  */}
        <Show when={instance() !== undefined && !instance().is_initialized}>
          <Navigate href="/getting-started" />
        </Show>

        <Routes>
          <Route path="/getting-started" component={GettingStarted} />
        </Routes>
      </Suspense>
    </Layout >
  )
}

export default App
