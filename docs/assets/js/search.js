// Search functionality for AWS Research Wizard docs
(function() {
  let searchIndex;
  let searchData;

  // Initialize search when page loads
  document.addEventListener('DOMContentLoaded', function() {
    initializeSearch();
  });

  function initializeSearch() {
    // Load search data
    fetch('/search.json')
      .then(response => response.json())
      .then(data => {
        searchData = data;

        // Build search index
        searchIndex = lunr(function() {
          this.field('title', { boost: 10 });
          this.field('content', { boost: 5 });
          this.field('category');
          this.field('tags');
          this.ref('id');

          data.forEach(function(doc) {
            this.add(doc);
          }, this);
        });

        // Set up search input handler
        const searchInput = document.getElementById('search-input');
        const searchResults = document.getElementById('search-results');

        if (searchInput) {
          searchInput.addEventListener('input', debounce(handleSearch, 300));

          // Hide results when clicking outside
          document.addEventListener('click', function(e) {
            if (!e.target.closest('.search-container')) {
              searchResults.style.display = 'none';
            }
          });
        }
      })
      .catch(error => {
        console.error('Error loading search data:', error);
      });
  }

  function handleSearch(event) {
    const query = event.target.value.trim();
    const searchResults = document.getElementById('search-results');

    if (query.length < 2) {
      searchResults.style.display = 'none';
      return;
    }

    try {
      // Perform search
      const results = searchIndex.search(query);

      // Display results
      displaySearchResults(results.slice(0, 10)); // Show top 10 results

    } catch (error) {
      console.error('Search error:', error);
      searchResults.style.display = 'none';
    }
  }

  function displaySearchResults(results) {
    const searchResults = document.getElementById('search-results');

    if (results.length === 0) {
      searchResults.innerHTML = '<div class="search-result">No results found</div>';
      searchResults.style.display = 'block';
      return;
    }

    const html = results.map(result => {
      const doc = searchData.find(d => d.id === result.ref);
      if (!doc) return '';

      return `
        <div class="search-result" onclick="window.location.href='${doc.url}'">
          <div class="search-result-title">${doc.title}</div>
          <div class="search-result-snippet">${truncate(doc.content, 120)}</div>
        </div>
      `;
    }).join('');

    searchResults.innerHTML = html;
    searchResults.style.display = 'block';
  }

  function truncate(text, length) {
    if (text.length <= length) return text;
    return text.substring(0, length) + '...';
  }

  function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout);
        func(...args);
      };
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
    };
  }
})();
