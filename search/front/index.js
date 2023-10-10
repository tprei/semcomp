document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.getElementById("searchInput");
    let delayTimer;

    searchInput.addEventListener("input", () => {
        clearTimeout(delayTimer); // Clear the previous timer
        delayTimer = setTimeout(() => {
            const query = searchInput.value;
            if (query) {
                queryAPI(query);
            } else {
                // Clear results when search input is empty
                updateResults([]);
            }
        }, 500); // Delay for 500 milliseconds before making the query
    });

    function queryAPI(query) {
        // Replace this with your actual API request
        // The response should contain data in the format:
        // [{ artist: "", title: "", img: "url_to_image" }]

        // Mock data for testing
        const data = [
            { artist: "Artist 1", title: "Song 1", img: "https://akamai.sscdn.co/letras/desktop/static/img/ic_placeholder_artist.svg" },
            { artist: "Artist 2", title: "Song 2", img: "https://akamai.sscdn.co/letras/desktop/static/img/ic_placeholder_artist.svg" },
        ];

        updateResults(data);
    }

    function updateResults(data) {
        const list = document.getElementById("results");

        // Clear previous results
        list.innerHTML = "";

        // Populate the results
        data.slice(0, 10).forEach(result => {
            const listItem = document.createElement("li");
            const avatar = document.createElement("img");
            avatar.src = result.img;
            avatar.width = 50;
            avatar.height = 50;
            listItem.appendChild(avatar);

            const text = document.createElement("p");
            text.innerText = `${result.artist} - ${result.title}`;
            listItem.appendChild(text);

            list.appendChild(listItem);
        });

        // Determine whether to slide down or slide up based on results
        // Set the max-height to the actual scroll height for animation
        list.style.maxHeight = list.scrollHeight + "px";
    }
});
