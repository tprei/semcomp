document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.getElementById("searchInput");
    const resultsList = document.getElementById("results");

    const prevPageButton = document.getElementById("prevPage");
    const nextPageButton = document.getElementById("nextPage");
    const currentPageText = document.getElementById("currentPage");

    let data = null;

    let pageSize = 5;
    let currentPage = 0; // Track the current page

    let query = "";

    function updatePagination() {
        // Enable or disable pagination buttons based on the current page
        prevPageButton.disabled = currentPage === 0;
        nextPageButton.disabled = (currentPage + 1) * pageSize >= data.length;

        // Update the current page indicator
        currentPageText.textContent = `Página ${currentPage + 1}`;

        updateResults();
    }

    prevPageButton.addEventListener("click", () => {
        if (currentPage > 0) {
            currentPage--;
            updatePagination();
        }
    });

    nextPageButton.addEventListener("click", () => {
        if ((currentPage + 1) * pageSize < data.length) {
            currentPage++;
            updatePagination();
        }
    });

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

            query = searchInput.value;
            if (query) {
                queryAPI(query);
            } else {
                // Clear results when search input is empty
                updateResults([]);
            }
        }, 500); // Delay for 500 milliseconds before making the query
    });

    async function queryAPI() {
        try {
            const response = await fetch(`http://localhost:8080/search?query=${query}`);
            if (response.ok) {
                data = await response.json();

                // Sort the data by views in descending order
                data.sort((a, b) => {
                    return parseInt(b.views.replaceAll('.', '')) - parseInt(a.views.replaceAll('.', ''));
                });

                currentPage = 0; // Reset current page to 0
                updatePagination();
            } else {
                data = null;
                // Handle the error if the response is not OK (e.g., network error)
                console.error("Error fetching data");
            }
        } catch (error) {
            // Handle any other errors (e.g., JSON parsing error)
            console.error("Error:", error);
        }
    }

    function updateResults() {
        const list = document.getElementById("results");

        // Clear previous results
        list.innerHTML = "";


        if (!data || data.length === 0) {
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

            const start = currentPage * pageSize;
            const end = Math.min(start + pageSize, data.length);

            // Populate the results for the current page
            for (let i = start; i < end; i++) {
                const result = data[i];
                const listItem = document.createElement("li");

                const divItem = document.createElement("div");
                divItem.className = "avatar"

                const avatar = document.createElement("img");
                avatar.src = result.img;
                avatar.width = 50;
                avatar.height = 50;

                const text = document.createElement("p");
                text.innerHTML = `<b>${result.artist}</b> - ${result.title}`;

                divItem.appendChild(avatar);
                divItem.appendChild(text);

                listItem.appendChild(divItem);

                const viewsContainer = document.createElement("div");
                viewsContainer.style.display = "flex";
                viewsContainer.style.flexDirection = "row";
                viewsContainer.style.alignItems = "center";
                viewsContainer.style.justifyContent = "flex-end";
                viewsContainer.style.gap = "10px";

                const viewsText = document.createElement("p");
                viewsText.innerText = result.views;
                viewsText.style.fontWeight = "bold";
                viewsContainer.appendChild(viewsText);

                const viewsIcon = document.createElement("img");
                viewsIcon.src = "assets/eye.svg";
                viewsIcon.style.height = '40px';
                viewsIcon.style.width = '40px';
                viewsContainer.appendChild(viewsIcon);

                listItem.appendChild(viewsContainer);

                // get lyrics
                let parts = "";
                for (let l = 0; l < result.lyrics.length; l++) {
                    for (let s = 0; s < result.lyrics[l].length; s++) {
                        let verse = result.lyrics[l][s];

                        console.log(verse, verse.toLowerCase().replace(/,!\.;\'/g, ''))
                        if (verse.toLowerCase().replace(/,!\.;\'/g, '').replace(/&#39;/g, "'").includes(query)) {
                            parts += verse.replace(/&#39;/g, "'") + "\n"
                        }
                    }
                }

                parts = parts.toLowerCase().replaceAll(query, `<b>${query}</b>`)

                const lyricsContainer = document.querySelector(".lyrics-container");
                const lyrics = document.querySelector(".lyrics");

                listItem.addEventListener("mouseenter", () => {
                    lyrics.style.display = "block";

                    if (query)
                        lyrics.innerHTML = parts;
                });

                listItem.addEventListener("mousemove", (event) => {
                    const x = event.clientX - lyricsContainer.getBoundingClientRect().left;
                    const y = event.clientY - lyricsContainer.getBoundingClientRect().top;

                    lyrics.style.left = x + 15 + "px"; // Add an offset to avoid cursor overlap
                    lyrics.style.top = y + 15 + "px";  // Add an offset to avoid cursor overlap
                });

                listItem.addEventListener("mouseleave", () => {
                    lyrics.style.display = "none";
                });

                list.appendChild(listItem);
            }

        }

        // Set the max-height to the actual scroll height for animation
        list.style.maxHeight = list.scrollHeight + "px";
    }
});
