<script>
  export let status

  function formatDate(dateString) {
    if (!dateString || dateString === '0001-01-01T00:00:00Z') return 'Never'
    const date = new Date(dateString)
    return date.toLocaleString()
  }

  $: isIndexing = status?.is_indexing || false
  $: lastIndexed = status?.last_indexed
  $: nextScheduled = status?.next_scheduled
  $: lastError = status?.last_error
  $: indexedPaths = status?.indexed_paths || []
</script>

<div class="space-y-4">
  <h2 class="text-xl font-semibold text-gray-800 border-b pb-2">Status</h2>

  <!-- Indexing Status -->
  <div class="flex items-center space-x-3">
    <div class="flex-shrink-0">
      {#if isIndexing}
        <div class="h-3 w-3 bg-yellow-500 rounded-full animate-pulse"></div>
      {:else}
        <div class="h-3 w-3 bg-green-500 rounded-full"></div>
      {/if}
    </div>
    <div>
      <p class="text-sm font-medium text-gray-700">
        {isIndexing ? 'Indexing in progress...' : 'Idle'}
      </p>
    </div>
  </div>

  <!-- Last Indexed -->
  <div class="bg-gray-50 rounded p-3">
    <p class="text-sm text-gray-600">Last Indexed</p>
    <p class="text-lg font-medium text-gray-800">{formatDate(lastIndexed)}</p>
  </div>

  <!-- Next Scheduled -->
  {#if nextScheduled && nextScheduled !== '0001-01-01T00:00:00Z'}
    <div class="bg-gray-50 rounded p-3">
      <p class="text-sm text-gray-600">Next Scheduled</p>
      <p class="text-lg font-medium text-gray-800">{formatDate(nextScheduled)}</p>
    </div>
  {/if}

  <!-- Indexed Paths -->
  {#if indexedPaths.length > 0}
    <div class="bg-gray-50 rounded p-3">
      <p class="text-sm text-gray-600 mb-2">Indexed Paths</p>
      <ul class="space-y-1">
        {#each indexedPaths as path}
          <li class="text-sm font-mono text-gray-700 truncate" title={path}>
            📁 {path}
          </li>
        {/each}
      </ul>
    </div>
  {/if}

  <!-- Error Display -->
  {#if lastError}
    <div class="bg-red-50 border border-red-200 rounded p-3">
      <p class="text-sm font-medium text-red-800 mb-1">Last Error</p>
      <p class="text-xs text-red-700 font-mono">{lastError}</p>
    </div>
  {/if}
</div>
