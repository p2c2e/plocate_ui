<script>
  import { createEventDispatcher } from 'svelte'

  export let status

  const dispatch = createEventDispatcher()

  let loading = false
  let indexLoading = {}
  let newIndexName = ''
  let newIndexPath = ''
  let addingIndex = false

  async function startIndexing(indexName = null) {
    if (indexName) {
      indexLoading[indexName] = true
    } else {
      loading = true
    }

    try {
      const url = indexName ? `/api/control/start/${indexName}` : '/api/control/start'
      const response = await fetch(url, { method: 'POST' })
      if (response.ok) {
        dispatch('statuschange')
      } else {
        const data = await response.json()
        alert(`Failed to start indexing: ${data.error}`)
      }
    } catch (error) {
      alert(`Error: ${error.message}`)
    } finally {
      if (indexName) {
        indexLoading[indexName] = false
      } else {
        loading = false
      }
    }
  }

  async function stopIndexing(indexName = null) {
    if (indexName) {
      indexLoading[indexName] = true
    } else {
      loading = true
    }

    try {
      const url = indexName ? `/api/control/stop/${indexName}` : '/api/control/stop'
      const response = await fetch(url, { method: 'POST' })
      if (response.ok) {
        dispatch('statuschange')
      } else {
        const data = await response.json()
        alert(`Failed to stop indexing: ${data.error}`)
      }
    } catch (error) {
      alert(`Error: ${error.message}`)
    } finally {
      if (indexName) {
        indexLoading[indexName] = false
      } else {
        loading = false
      }
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

  async function addIndex() {
    if (!newIndexName.trim() || !newIndexPath.trim()) return
    addingIndex = true

    try {
      const response = await fetch('/api/indices', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: newIndexName.trim(),
          index_paths: [newIndexPath.trim()]
        })
      })
      if (response.ok) {
        newIndexName = ''
        newIndexPath = ''
        dispatch('statuschange')
      } else {
        const data = await response.json()
        alert(`Failed to add index: ${data.error}`)
      }
    } catch (error) {
      alert(`Error: ${error.message}`)
    } finally {
      addingIndex = false
    }
  }

  async function removeIndex(indexName) {
    if (!confirm(`Remove index "${indexName}"? This will stop indexing and remove the configuration.`)) return
    indexLoading[indexName] = true

    try {
      const response = await fetch(`/api/indices/${indexName}`, { method: 'DELETE' })
      if (response.ok) {
        dispatch('statuschange')
      } else {
        const data = await response.json()
        alert(`Failed to remove index: ${data.error}`)
      }
    } catch (error) {
      alert(`Error: ${error.message}`)
    } finally {
      indexLoading[indexName] = false
    }
  }

  function formatDate(dateStr) {
    if (!dateStr || dateStr === '0001-01-01T00:00:00Z') return 'Never'
    const date = new Date(dateStr)
    return date.toLocaleString()
  }

  $: indices = status?.indices || []
  $: hasSchedule = status?.next_scheduled && status.next_scheduled !== '0001-01-01T00:00:00Z'
  $: anyIndexing = indices.some(idx => idx.is_indexing)
</script>

<div class="space-y-4">
  <h2 class="text-xl font-semibold text-gray-800 border-b pb-2">Controls</h2>

  <!-- Global Control -->
  <div class="space-y-2">
    <p class="text-sm font-medium text-gray-700">All Indices</p>
    <div class="flex space-x-2">
      <button
        on:click={() => startIndexing()}
        disabled={loading || anyIndexing}
        class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium text-sm"
      >
        {#if anyIndexing}
          Indexing...
        {:else}
          Start All
        {/if}
      </button>

      <button
        on:click={() => stopIndexing()}
        disabled={loading || !anyIndexing}
        class="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium text-sm"
      >
        Stop All
      </button>
    </div>
  </div>

  <!-- Add Index -->
  <div class="space-y-2">
    <p class="text-sm font-medium text-gray-700">Add Index</p>
    <div class="bg-gray-50 border border-gray-200 rounded-lg p-3 space-y-2">
      <input
        type="text"
        bind:value={newIndexName}
        placeholder="Index name (e.g. documents)"
        class="w-full px-3 py-2 border border-gray-300 rounded text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
      />
      <input
        type="text"
        bind:value={newIndexPath}
        placeholder="Folder path (e.g. /mnt/Documents)"
        class="w-full px-3 py-2 border border-gray-300 rounded text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
      />
      <button
        on:click={addIndex}
        disabled={addingIndex || !newIndexName.trim() || !newIndexPath.trim()}
        class="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium text-sm"
      >
        {#if addingIndex}
          Adding...
        {:else}
          Add Index
        {/if}
      </button>
    </div>
  </div>

  <!-- Per-Index Control -->
  {#if indices.length > 0}
    <div class="space-y-3">
      <p class="text-sm font-medium text-gray-700">Individual Indices</p>
      {#each indices as index}
        <div class="bg-gray-50 border border-gray-200 rounded-lg p-3">
          <div class="flex items-start justify-between mb-2">
            <div class="flex-1">
              <div class="flex items-center space-x-2">
                <h3 class="font-medium text-gray-800">{index.name}</h3>
                {#if !index.enabled}
                  <span class="px-2 py-0.5 text-xs bg-gray-300 text-gray-700 rounded">Disabled</span>
                {/if}
                {#if index.is_indexing}
                  <span class="px-2 py-0.5 text-xs bg-blue-500 text-white rounded animate-pulse">Indexing</span>
                {/if}
              </div>
              <p class="text-xs text-gray-500 mt-1">
                Paths: {index.indexed_paths.join(', ')}
              </p>
              <p class="text-xs text-gray-500">
                Last indexed: {formatDate(index.last_indexed)}
              </p>
              {#if index.last_error}
                <p class="text-xs text-red-600 mt-1">
                  Error: {index.last_error}
                </p>
              {/if}
            </div>
            <button
              on:click={() => removeIndex(index.name)}
              disabled={indexLoading[index.name]}
              class="flex-shrink-0 ml-2 px-2 py-1 text-xs text-red-600 hover:bg-red-100 rounded transition-colors"
              title="Remove index"
            >
              Remove
            </button>
          </div>
          <div class="flex space-x-2">
            <button
              on:click={() => startIndexing(index.name)}
              disabled={indexLoading[index.name] || index.is_indexing || !index.enabled}
              class="flex-1 px-3 py-1.5 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors text-xs font-medium"
            >
              Start
            </button>
            <button
              on:click={() => stopIndexing(index.name)}
              disabled={indexLoading[index.name] || !index.is_indexing}
              class="flex-1 px-3 py-1.5 bg-red-600 text-white rounded hover:bg-red-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors text-xs font-medium"
            >
              Stop
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}

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
        ✓ Next run: {formatDate(status.next_scheduled)}
      </p>
    {/if}
  </div>

  <!-- Info -->
  <div class="bg-blue-50 border border-blue-200 rounded p-3 mt-4">
    <p class="text-xs text-blue-800">
      <strong>Tip:</strong> You can index all enabled indices at once or index them individually.
      The scheduler will automatically index all enabled indices.
    </p>
  </div>
</div>
