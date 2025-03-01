document.addEventListener('DOMContentLoaded', function() {
    const form = document.querySelector('.action-form');
    const button = form.querySelector('.action-form__button');
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

    form.addEventListener('submit', (event) => {
        event.preventDefault();

        const formData = new FormData(form);
        const webPageUrl = formData.get('url');

        const payload = {webPageUrl};

        console.log('Payload:', payload);

        button.disabled = true;

        fetch('/api/analyze', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(payload),
        })
        .then(response => response.json())
        .then((data) => {
            console.log('Success:', data);
            // const result = document.querySelector('.result');
            // result.innerHTML = `<a href="${data.url}" target="_blank">${data.url}</a>`;
        })
        .catch(error => {
            console.error('Error:', error);
        })
        .finally(() => {
            button.disabled = false;
        });
    });
});
