document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.getElementById("searchInput");
    const searchButton = document.getElementById("searchButton");
    const resultsList = document.getElementById("results");

    searchButton.addEventListener("click", () => {
        const query = searchInput.value;
        if (query) {
            queryAPI(query);
        }
    });

    function queryAPI(query) {
        fetch(`http://localhost:8080/search?query=${query}`)
            .then(response => response.json())
            .then(data => {
                updateResults(data);
            })
            .catch(error => {
                console.error(error);
            });
    }

    function updateResults(data) {
        resultsList.innerHTML = ""; // Clear the previous results
        data.slice(0, 10).forEach(result => {
            const listItem = document.createElement("li");
            listItem.innerHTML = `${result.artist} - ${result.title}`;
            resultsList.appendChild(listItem);
        });
    }
});
