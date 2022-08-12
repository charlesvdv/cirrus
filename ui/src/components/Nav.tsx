import { createSignal } from "solid-js"

import Logo from './Logo'

export default () => {
  const [showMenu, setShowMenu] = createSignal(false);

  return (
    <>
      <nav class="bg-white shadow">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16">
            <div class="flex">
              <div class="flex-shrink-0 flex items-center space-x-2">
                <Logo class="block h-10 w-auto fill-primary" />
                <div class="text-lg font-semibold text-primary">cirrus</div>
              </div>
              <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
                {/* TODO: uncomment this when menu is needed */}
                {/* <a href="#" class="border-primary text-gray-900 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"> Dashboard </a>
                    <a href="#" class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"> Team </a>
                    <a href="#" class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"> Projects </a>
                    <a href="#" class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"> Calendar </a> */}
              </div>
            </div>

            {/* TODO: uncomment this when menu is needed */}
            {/* <div class="-mr-2 flex items-center sm:hidden">
              <button type="button" onClick={() => setShowMenu(!showMenu())} class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary" aria-controls="mobile-menu" aria-expanded="false">
                <span class="sr-only">Open main menu</span>
                <svg classList={{ hidden: !showMenu() }} class="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
                </svg>
                <svg classList={{ hidden: showMenu() }} class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div> */}
          </div>
        </div>

        {/* TODO: uncomment this when menu is needed */}
        {/* <div classList={{ hidden: showMenu() }} class="sm:hidden" id="mobile-menu">
          <div class="pt-2 pb-3 space-y-1">
            <a href="#" class="bg-indigo-50 border-primary text-indigo-700 block pl-3 pr-4 py-2 border-l-4 text-base font-medium">Dashboard</a>
            <a href="#" class="border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700 block pl-3 pr-4 py-2 border-l-4 text-base font-medium">Team</a>
            <a href="#" class="border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700 block pl-3 pr-4 py-2 border-l-4 text-base font-medium">Projects</a>
            <a href="#" class="border-transparent text-gray-500 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-700 block pl-3 pr-4 py-2 border-l-4 text-base font-medium">Calendar</a>
          </div>
        </div> */}
      </nav>
    </>
  )
}