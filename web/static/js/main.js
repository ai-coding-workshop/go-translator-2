document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('translation-form');
    const translateBtn = document.getElementById('translate-btn');
    const btnText = document.querySelector('.btn-text');
    const btnLoading = document.querySelector('.btn-loading');
    const resultContainer = document.getElementById('result-container');
    const originalText = document.getElementById('original-text');
    const modelUsed = document.getElementById('model-used');
    const translationResult = document.getElementById('translation-result');
    const newTranslationBtn = document.getElementById('new-translation-btn');

    // Handle form submission
    form.addEventListener('submit', function(e) {
        e.preventDefault();

        // Get form data
        const text = document.getElementById('text').value;
        const model = document.getElementById('model').value;

        // Show loading state
        btnText.style.display = 'none';
        btnLoading.style.display = 'inline';
        translateBtn.disabled = true;

        // Prepare data for API call
        const data = {
            text: text,
            model: model
        };

        // Call API for translation
        fetch('/api/translate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        })
        .then(response => response.json())
        .then(data => {
            // Display result
            originalText.textContent = text;
            modelUsed.textContent = getModelName(model);
            translationResult.textContent = data.translation || 'Translation will appear here';

            // Hide form and show result
            form.style.display = 'none';
            resultContainer.style.display = 'block';
        })
        .catch(error => {
            console.error('Error:', error);
            translationResult.textContent = 'Error occurred during translation. Please try again.';
            resultContainer.style.display = 'block';
        })
        .finally(() => {
            // Reset button state
            btnText.style.display = 'inline';
            btnLoading.style.display = 'none';
            translateBtn.disabled = false;
        });
    });

    // Handle "Translate Another" button
    newTranslationBtn.addEventListener('click', function() {
        // Hide result and show form
        resultContainer.style.display = 'none';
        form.style.display = 'block';

        // Clear form
        document.getElementById('text').value = '';
        document.getElementById('model').selectedIndex = 0;
    });

    // Helper function to get model name
    function getModelName(modelKey) {
        const models = {
            'gpt-3.5': 'GPT-3.5 Turbo',
            'gpt-4': 'GPT-4',
            'claude': 'Claude',
            'llama': 'Llama 2'
        };
        return models[modelKey] || modelKey;
    }
});
