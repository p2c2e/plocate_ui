<script>
  export let status

  function formatDate(dateString) {
    if (!dateString || dateString === '0001-01-01T00:00:00Z') return 'Never'
    const date = new Date(dateString)
    return date.toLocaleString()
  }

  $: indices = status?.indices || []
  $: nextScheduled = status?.next_scheduled
  $: anyIndexing = indices.some(idx => idx.is_indexing)
  $: totalPaths = indices.reduce((acc, idx) => acc + (idx.indexed_paths?.length || 0), 0)
</script>

<div class="space-y-4">
  <h2 class="text-xl font-semibold text-gray-800 border-b pb-2">Status</h2>

  <!-- Overall Status -->
  <div class="flex items-center space-x-3">
    <div class="flex-shrink-0">
      {#if anyIndexing}
        <div class="h-3 w-3 bg-yellow-500 rounded-full animate-pulse"></div>
      {:else}
        <div class="h-3 w-3 bg-green-500 rounded-full"></div>
      {/if}
    </div>
    <div>
      <p class="text-sm font-medium text-gray-700">
        {anyIndexing ? 'Indexing in progress...' : 'All indices idle'}
      </p>
    </div>
  </div>

  <!-- Summary Stats -->
  <div class="grid grid-cols-2 gap-3">
    <div class="bg-gray-50 rounded p-3">
      <p class="text-sm text-gray-600">Total Indices</p>
      <p class="text-lg font-medium text-gray-800">{indices.length}</p>
    </div>
    <div class="bg-gray-50 rounded p-3">
      <p class="text-sm text-gray-600">Total Paths</p>
      <p class="text-lg font-medium text-gray-800">{totalPaths}</p>
    </div>
  </div>

  <!-- Next Scheduled -->
  {#if nextScheduled && nextScheduled !== '0001-01-01T00:00:00Z'}
    <div class="bg-gray-50 rounded p-3">
      <p class="text-sm text-gray-600">Next Scheduled Run</p>
      <p class="text-lg font-medium text-gray-800">{formatDate(nextScheduled)}</p>
    </div>
  {/if}

  <!-- Per-Index Status Summary -->
  {#if indices.length > 0}
    <div class="bg-gray-50 rounded p-3">
      <p class="text-sm text-gray-600 mb-2">Indices Overview</p>
      <div class="space-y-2">
        {#each indices as index}
          <div class="bg-white rounded p-2 border border-gray-200">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <span class="font-medium text-gray-800 text-sm">{index.name}</span>
                {#if index.is_indexing}
                  <span class="px-2 py-0.5 text-xs bg-yellow-500 text-white rounded">Indexing</span>
                {:else if !index.enabled}
                  <span class="px-2 py-0.5 text-xs bg-gray-300 text-gray-700 rounded">Disabled</span>
                {:else}
                  <span class="px-2 py-0.5 text-xs bg-green-500 text-white rounded">Ready</span>
                {/if}
              </div>
              <span class="text-xs text-gray-500">{index.indexed_paths?.length || 0} paths</span>
            </div>
            {#if index.last_indexed && index.last_indexed !== '0001-01-01T00:00:00Z'}
              <p class="text-xs text-gray-500 mt-1">
                Last: {formatDate(index.last_indexed)}
              </p>
            {/if}
            {#if index.last_error}
              <p class="text-xs text-red-600 mt-1 truncate" title={index.last_error}>
                ⚠️ {index.last_error}
              </p>
            {/if}
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>
