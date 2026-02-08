<script>
  import { createEventDispatcher } from 'svelte'

  export let status

  const dispatch = createEventDispatcher()

  let loading = false

  async function startIndexing() {
    loading = true
    try {
      const response = await fetch('/api/control/start', { method: 'POST' })
      if (response.ok) {
        dispatch('statuschange')
      } else {
        const data = await response.json()
        alert(`Failed to start indexing: ${data.error}`)
      }
    } catch (error) {
      alert(`Error: ${error.message}`)
    } finally {
      loading = false
    }
  }

  async function stopIndexing() {
    loading = true
    try {
      const response = await fetch('/api/control/stop', { method: 'POST' })
      if (response.ok) {
        dispatch('statuschange')
      } else {
        const data = await response.json()
        alert(`Failed to stop indexing: ${data.error}`)
      }
    } catch (error) {
      alert(`Error: ${error.message}`)
    } finally {
      loading = false
    }
  }

  async function enableScheduler() {
    loading = true
    try {
      const response = await fetch('/api/control/scheduler/enable', { method: 'POST' })
      if (response.ok) {
        dispatch('statuschange')
      }
    } catch (error) {
      alert(`Error: ${error.message}`)
    } finally {
      loading = false
    }
  }

  async function disableScheduler() {
    loading = true
    try {
      const response = await fetch('/api/control/scheduler/disable', { method: 'POST' })
      if (response.ok) {
        dispatch('statuschange')
      }
    } catch (error) {
      alert(`Error: ${error.message}`)
    } finally {
      loading = false
    }
  }

  $: isIndexing = status?.is_indexing || false
  $: hasSchedule = status?.next_scheduled && status.next_scheduled !== '0001-01-01T00:00:00Z'
</script>

<div class="space-y-4">
  <h2 class="text-xl font-semibold text-gray-800 border-b pb-2">Controls</h2>

  <!-- Manual Control -->
  <div class="space-y-2">
    <p class="text-sm font-medium text-gray-700">Manual Indexing</p>
    <div class="flex space-x-2">
      <button
        on:click={startIndexing}
        disabled={loading || isIndexing}
        class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium text-sm"
      >
        {#if isIndexing}
          Indexing...
        {:else}
          Start Index Now
        {/if}
      </button>

      <button
        on:click={stopIndexing}
        disabled={loading || !isIndexing}
        class="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium text-sm"
      >
        Stop Indexing
      </button>
    </div>
  </div>

  <!-- Scheduler Control -->
  <div class="space-y-2">
    <p class="text-sm font-medium text-gray-700">Automatic Scheduler</p>
    <div class="flex space-x-2">
      <button
        on:click={enableScheduler}
        disabled={loading || hasSchedule}
        class="flex-1 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium text-sm"
      >
        Enable Scheduler
      </button>

      <button
        on:click={disableScheduler}
        disabled={loading || !hasSchedule}
        class="flex-1 px-4 py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium text-sm"
      >
        Disable Scheduler
      </button>
    </div>
    {#if hasSchedule}
      <p class="text-xs text-gray-600 mt-1">
        ✓ Scheduler is active
      </p>
    {/if}
  </div>

  <!-- Info -->
  <div class="bg-blue-50 border border-blue-200 rounded p-3 mt-4">
    <p class="text-xs text-blue-800">
      <strong>Tip:</strong> The scheduler runs automatically at configured intervals.
      You can also trigger manual indexing at any time.
    </p>
  </div>
</div>
