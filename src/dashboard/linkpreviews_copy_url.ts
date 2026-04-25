/** Initializes the link preview URL generator UI by attaching event listeners to buttons. */
export function initLinkPreviewsCopyUrl(): void {
	// Attach event listeners to buttons
	const generateBtn = document.getElementById('generate-link-preview-btn');
	const copyBtn = document.getElementById('copy-link-preview-btn');

	if (generateBtn) {
		generateBtn.addEventListener('click', generateLinkPreviewUrl);
	}

	if (copyBtn) {
		copyBtn.addEventListener('click', copyToClipboard);
	}
}

/**
 * Generates a link preview URL for the given input URL and displays it for copying.
 *
 * Uses HTML5 form validation to ensure the input field is filled and contains a valid URL format.
 * If valid, constructs the full `/link-previews/v1` URL with the input encoded as a query parameter
 * and displays it in the output field.
 */
export function generateLinkPreviewUrl(): void {
	const inputElement = document.getElementById('generate-link-url') as HTMLInputElement;
	const form = inputElement.closest('form') as HTMLFormElement;

	// Use HTML5 validation
	if (!form.checkValidity()) {
		inputElement.reportValidity();
		return;
	}

	const url = inputElement.value.trim();

	try {
		new URL(url);
	} catch {
		alert('Please enter a valid URL');
		return;
	}

	const baseUrl = window.location.origin;
	const linkPreviewUrl = `${baseUrl}/link-previews/v1?url=${encodeURIComponent(url)}`;
	(document.getElementById('generated-link-url') as HTMLInputElement).value = linkPreviewUrl;
	document.getElementById('generated-link-output')!.classList.remove('hidden');
}

/**
 * Copies the generated link preview URL to the clipboard and provides visual feedback.
 * Temporarily changes the button text to "Copied!" for 2 seconds to provide user feedback,
 * then restores the original button text.
 */
export function copyToClipboard(event: Event): void {
	const linkUrl = document.getElementById('generated-link-url') as HTMLInputElement;
	linkUrl.select();
	document.execCommand('copy');
	const btn = event.target as HTMLButtonElement;
	const originalText = btn.textContent;
	btn.textContent = 'Copied!';
	setTimeout(() => {
		btn.textContent = originalText;
	}, 2000);
}
