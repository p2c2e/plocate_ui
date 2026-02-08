<script>
  import Search from './lib/Search.svelte'
  import Status from './lib/Status.svelte'
  import Controls from './lib/Controls.svelte'
  import { onMount } from 'svelte'

  let status = null
  let statusInterval = null

  async function fetchStatus() {
    try {
      const response = await fetch('/api/status')
      status = await response.json()
    } catch (error) {
      console.error('Failed to fetch status:', error)
    }
  }

  onMount(() => {
    fetchStatus()
    statusInterval = setInterval(fetchStatus, 5000) // Update every 5 seconds

    return () => {
      if (statusInterval) clearInterval(statusInterval)
    }
  })

  function handleStatusChange() {
    fetchStatus()
  }
</script>

<main class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
  <div class="container mx-auto px-4 py-8 max-w-7xl">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-4xl font-bold text-gray-800 mb-2">
        <span class="text-blue-600">Plocate</span> File Search
      </h1>
      <p class="text-gray-600">Fast file location service for your Unraid server</p>
    </div>

    <!-- Status and Controls Card -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-6">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Status {status} />
        <Controls {status} on:statuschange={handleStatusChange} />
      </div>
    </div>

    <!-- Search Card -->
    <div class="bg-white rounded-lg shadow-md p-6">
      <Search />
    </div>
  </div>
</main>
