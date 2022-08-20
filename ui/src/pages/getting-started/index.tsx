import { useNavigate } from "@solidjs/router";

import { useState } from "../../lib/state"
import Logo from "../../components/Logo"

export default () => {
  const navigate = useNavigate()
  const { client } = useState()

  const onSubmit = async (e) => {
    e.preventDefault()
    const formData = Object.fromEntries(new FormData(e.target))

    await client.initInstance({
      admin: {
        name: formData.username as string,
        password: formData.password as string,
      }
    })

    navigate("/login")
  }

  return (
    <div class="min-h-full flex flex-col items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full my-16">
        <Logo class="mx-auto h-24 fill-primary" />

        <h1 class="mt-6 text-center text-3xl font-extrabold text-gray-900">Welcome to cirrus!</h1>
        <p class="text-center text-gray-500 text-xl">Fast and simple personal cloud solution</p>
      </div>

      <h2 class="text-gray-900 text-center font-semibold -my-6">Create an admin account</h2>

      <div class="max-w-md w-full space-y-8">
        <form class="mt-8 space-y-6" onSubmit={onSubmit} >
          <div class="rounded-md shadow-sm -space-y-px">
            <div>
              <label for="username" class="sr-only">Username</label>
              <input id="username" name="username" type="text" required class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-primary focus:border-primary focus:z-10 sm:text-sm" placeholder="Username" />
            </div>
            <div>
              <label for="password" class="sr-only">Password</label>
              <input id="password" name="password" type="password" autocomplete="current-password" required class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-primary focus:border-primary focus:z-10 sm:text-sm" placeholder="Password" />
            </div>
          </div>

          <div>
            <button type="submit" class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary hover:bg-secondary focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary">
              Get started!
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}