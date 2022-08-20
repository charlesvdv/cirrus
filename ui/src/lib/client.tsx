const API_URL = "/api"

export class Client {
  constructor() {

  }

  async getInstance(): Promise<Instance> {
    return this.get("/instance")
  }

  async initInstance(req: InitInstance): Promise<void> {
    return this.post("/instance/init", req)
  }

  private async get<R>(url: string): Promise<R> {
    return this.send<void, R>("get", url, undefined)
  }

  private async post<T, R>(url: string, data: T): Promise<R> {
    return this.send<T, R>("post", url, data)
  }

  private async send<T, R>(method: string, url: string, data: T): Promise<R> {
    const headers = {},
      opts = { method, headers, body: null }

    if (data !== undefined) {
      headers["Content-Type"] = "application/json"
      opts.body = JSON.stringify(data)
    }

    try {
      const response = await fetch(API_URL + url, opts)
      return await response.json()
    } catch (err) {
      // TODO
    }
  }
}

export interface Instance {
  is_initialized: boolean,
}

export interface InitInstance {
  admin: {
    name: string,
    password: string,
  },
}