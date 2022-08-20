/* @refresh reload */
import { render } from 'solid-js/web'
import { Router } from "@solidjs/router";

import './assets/index.css'
import '@fontsource/inter/variable-full.css'

import App from './App'
import { StateProvider } from './lib/state';

render(
  () => (
    <Router>
      <StateProvider>
        <App />
      </StateProvider>
    </Router>
  ),
  document.getElementById('root') as HTMLElement
)
