document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.getElementById("searchInput");
    const resultsList = document.getElementById("results");

    let delayTimer;
    let errorItem = null; // Track the error message item

    searchInput.addEventListener("input", () => {
        clearTimeout(delayTimer); // Clear the previous timer
        delayTimer = setTimeout(() => {
            // Clear results and error message when search input is empty
            if (errorItem) {
                resultsList.removeChild(errorItem);
                errorItem = null;
            }

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
        if (!query) {
            const resultsList = document.getElementById("results");
            resultsList.replaceChildren();
            errorItem = null;
        }
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

        if (data.length === 0) {
            // Create a special list item for "No results were found"
            const noResultsItem = document.createElement("li");
            noResultsItem.innerText = "Puts, não achei essa! 😞";
            noResultsItem.style.backgroundColor = "#B91D82"; // Set the background color
            noResultsItem.style.color = "#FFFFFF"; // Set the text color
            noResultsItem.style.fontFamily = "Roboto";
            noResultsItem.style.fontWeight = "bold";
            noResultsItem.style.padding = "0.5rem";
            noResultsItem.style.borderRadius = "5px";

            list.appendChild(noResultsItem);

            // Store the error message item
            errorItem = noResultsItem;
        } else {
            // Clear the error message when results are found
            if (errorItem) {
                list.removeChild(errorItem);
                errorItem = null;
            }

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
        }

        // Set the max-height to the actual scroll height for animation
        list.style.maxHeight = list.scrollHeight + "px";
    }
});
