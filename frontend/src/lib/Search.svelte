<script>
  export let status

  let query = ''
  let results = []
  let loading = false
  let searchTime = 0
  let hasSearched = false
  let selectedIndices = []

  $: availableIndices = (status?.indices || []).map(idx => idx.name)

  // Auto-select new indices and remove stale selections
  $: {
    const newIndices = availableIndices.filter(name => !selectedIndices.includes(name))
    if (newIndices.length > 0) {
      selectedIndices = [...selectedIndices, ...newIndices]
    }
    selectedIndices = selectedIndices.filter(name => availableIndices.includes(name))
  }

  function toggleIndex(indexName) {
    if (selectedIndices.includes(indexName)) {
      selectedIndices = selectedIndices.filter(i => i !== indexName)
    } else {
      selectedIndices = [...selectedIndices, indexName]
    }
  }

  function selectAllIndices() {
    selectedIndices = [...availableIndices]
  }

  function deselectAllIndices() {
    selectedIndices = []
  }

  async function search() {
    if (!query.trim()) return

    loading = true
    hasSearched = true
    const startTime = performance.now()

    try {
      const requestBody = {
        query: query,
        limit: 500,
        indices: selectedIndices
      }

      const response = await fetch('/api/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestBody)
      })
      const data = await response.json()

      if (response.ok) {
        results = data.results || []
        searchTime = Math.round(performance.now() - startTime)
      } else {
        alert(`Search failed: ${data.error}`)
        results = []
      }
    } catch (error) {
      alert(`Error: ${error.message}`)
      results = []
    } finally {
      loading = false
    }
  }

  function handleKeyPress(event) {
    if (event.key === 'Enter') {
      search()
    }
  }

  function highlightMatch(path) {
    const parts = path.split('/')
    const filename = parts[parts.length - 1]
    const directory = parts.slice(0, -1).join('/')

    return {
      directory: directory || '/',
      filename
    }
  }
</script>

<div class="space-y-4">
  <h2 class="text-2xl font-semibold text-gray-800 mb-4">Search Files</h2>

  <!-- Index Selection -->
  {#if availableIndices.length > 0}
    <div class="bg-gray-50 border border-gray-200 rounded-lg p-3">
      <div class="flex items-center justify-between mb-2">
        <p class="text-sm font-medium text-gray-700">Search in indices:</p>
        <div class="flex space-x-2">
          <button
            on:click={selectAllIndices}
            class="text-xs text-blue-600 hover:text-blue-800 font-medium"
          >
            Select All
          </button>
          <span class="text-gray-400">|</span>
          <button
            on:click={deselectAllIndices}
            class="text-xs text-blue-600 hover:text-blue-800 font-medium"
          >
            Clear
          </button>
        </div>
      </div>
      <div class="flex flex-wrap gap-2">
        {#each availableIndices as indexName}
          <label class="flex items-center space-x-2 px-3 py-1 bg-white border border-gray-300 rounded-md hover:bg-gray-100 cursor-pointer transition-colors">
            <input
              type="checkbox"
              checked={selectedIndices.includes(indexName)}
              on:change={() => toggleIndex(indexName)}
              class="form-checkbox h-4 w-4 text-blue-600 rounded"
            />
            <span class="text-sm text-gray-700">{indexName}</span>
          </label>
        {/each}
      </div>
      {#if selectedIndices.length === 0}
        <p class="text-xs text-orange-600 mt-2">⚠️ No indices selected. Please select at least one index to search.</p>
      {/if}
    </div>
  {/if}

  <!-- Search Input -->
  <div class="flex space-x-2">
    <input
      type="text"
      bind:value={query}
      on:keypress={handleKeyPress}
      placeholder="Enter filename or pattern..."
      class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none text-lg"
      disabled={loading}
    />
    <button
      on:click={search}
      disabled={loading || !query.trim() || selectedIndices.length === 0}
      class="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium text-lg"
    >
      {#if loading}
        Searching...
      {:else}
        Search
      {/if}
    </button>
  </div>

  <!-- Results Summary -->
  {#if hasSearched && !loading}
    <div class="flex items-center justify-between text-sm text-gray-600">
      <p>
        Found <strong class="text-gray-800">{results.length}</strong>
        {results.length === 1 ? 'result' : 'results'}
      </p>
      <p>
        Search completed in <strong class="text-gray-800">{searchTime}ms</strong>
      </p>
    </div>
  {/if}

  <!-- Results List -->
  {#if results.length > 0}
    <div class="bg-gray-50 rounded-lg border border-gray-200 max-h-[600px] overflow-y-auto">
      <div class="divide-y divide-gray-200">
        {#each results as result}
          {@const parts = highlightMatch(result)}
          <div class="p-3 hover:bg-blue-50 transition-colors">
            <div class="flex items-start space-x-2">
              <span class="text-gray-400 mt-1">📄</span>
              <div class="flex-1 min-w-0">
                <p class="text-sm text-gray-500 truncate" title={parts.directory}>
                  {parts.directory}
                </p>
                <p class="text-base font-medium text-gray-800 break-all">
                  {parts.filename}
                </p>
              </div>
              <button
                on:click={() => navigator.clipboard.writeText(result)}
                class="flex-shrink-0 px-2 py-1 text-xs text-blue-600 hover:bg-blue-100 rounded transition-colors"
                title="Copy path"
              >
                Copy
              </button>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {:else if hasSearched && !loading}
    <div class="text-center py-12 bg-gray-50 rounded-lg">
      <p class="text-gray-500 text-lg">No results found for "{query}"</p>
      <p class="text-gray-400 text-sm mt-2">Try a different search term or select more indices</p>
    </div>
  {/if}

  <!-- Loading State -->
  {#if loading}
    <div class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      <p class="text-gray-600 mt-4">Searching...</p>
    </div>
  {/if}

  <!-- Initial State -->
  {#if !hasSearched && !loading}
    <div class="text-center py-12 bg-gray-50 rounded-lg">
      <p class="text-gray-500 text-lg">Start typing to search for files</p>
      <p class="text-gray-400 text-sm mt-2">Searches are case-insensitive and match partial names</p>
    </div>
  {/if}
</div>
