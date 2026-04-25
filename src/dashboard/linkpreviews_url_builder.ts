/** Initializes the link preview URL generator UI by attaching event listeners to buttons. */
export function initLinkPreviewsUrlBuilder(): void {
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
 * If valid, constructs the full `/link-previews/v1` URL with the input encoded as a query parameter
 * and displays it in the output field.
 */
export function generateLinkPreviewUrl(): void {
	const urlInput = document.getElementById('generate-link-url') as HTMLInputElement;
	const form = urlInput.closest('form') as HTMLFormElement;
	if (!form.checkValidity()) {
		urlInput.reportValidity();
		return;
	}

	const url = urlInput.value.trim();

	const baseUrl = window.location.origin;
  const params = new URLSearchParams();
	params.set('url', url);

  const linkPreviewUrl = `${baseUrl}/link-previews/v1?${params.toString()}`;
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
