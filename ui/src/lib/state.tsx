import { createContext, useContext } from "solid-js"
import { createStore } from "solid-js/store"

import { Client } from "./client"

interface State {
  client: Client,
  state: any,
}

export const StateContext = createContext<State>();

export function StateProvider(props) {
  const [state, setState] = createStore({
    authToken: undefined,
  })

  const client = new Client();

  return (
    <StateContext.Provider value={{ state, client }}>
      {props.children}
    </StateContext.Provider>
  )
}

export function useState(): State {
  return useContext(StateContext)
}