document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('.action-form');
    const loader = document.querySelector('.loader');
    const result = document.querySelector('.result');
    const button = form.querySelector('.action-form__button');
    const errorDisplay = document.querySelector('.error-display');
    button.disabled = true;

    form.querySelector('.action-form__input').addEventListener('input', (event) => {
        // URL validation using regex
        const urlPattern = new RegExp('^(https?:\\/\\/)?'+ // protocol
            '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.?)+[a-z]{2,}|'+ // domain name
            '((\\d{1,3}\\.){3}\\d{1,3}))'+ // OR ip (v4) address
            '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*'+ // port and path
            '(\\?[;&a-z\\d%_.~+=-]*)?'+ // query string
            '(\\#[-a-z\\d_]*)?$','i'); // fragment locator

        const value = event.target.value;
        const valid = urlPattern.test(value);
        
        if (valid) {
            button.disabled = false;
        } else {
            button.disabled = true;
        }
    });

    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const formData = new FormData(form);
        const webPageUrl = formData.get('url');

        const payload = {webPageUrl};

        console.log('Payload:', payload);

        button.disabled = true;

        loader.style.display = 'block';
        result.style.display = 'none';
        errorDisplay.style.display = 'none';

        fetch('/api/analyze', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(payload),
        })
        .then(async (response) => {
            if (!response.ok) {
                const error = await response.json();
                document.querySelector("#error-message").textContent = error.error;
                document.querySelector("#error-status").textContent = `Status Code: ${error.statusCode || 400}`;
                errorDisplay.style.display = 'block';
            } else {
                const data = await response.json();

                console.log('Success:', data);
            
                // Update result container with response data
                document.querySelector("#title").textContent = data.title;
                document.querySelector("#htmlVersion").textContent = data.htmlVersion;
                document.querySelector("#has-login-form").textContent = data.hasLoginForm ? "Yes" : "No";
                document.querySelector("#external-links").textContent = data.externalLinks;
                document.querySelector("#internal-links").textContent = data.internalLinks;
                document.querySelector("#inacc-links").textContent = data.inaccessibleLinks;

                // Update headings
                const headingsList = document.querySelector("#headings");
                headingsList.innerHTML = ""; // Clear existing list
                for (const [tag, count] of Object.entries(data.headings)) {
                    const listItem = document.createElement("li");
                    listItem.textContent = `${tag}: ${count}`;
                    headingsList.appendChild(listItem);
                }

                // Show result container
                result.style.display = "block";
            }
        })
        .catch(error => {
            console.error('Error:', error);
            document.querySelector("#error-message").textContent = "Something went wrong";
            document.querySelector("#error-status").textContent = `Status Code: ${500}`;
            errorDisplay.style.display = 'block';
        })
        .finally(() => {
            button.disabled = false;
            loader.style.display = 'none';
        });
    });
});
