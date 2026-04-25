/** Initializes the ratings embed builder UI by attaching event listeners to the form. */
export function initRatingsEmbedBuilder(): void {
	const generateBtn = document.getElementById('generate-rating-iframe-btn');
  const copyBtn = document.getElementById('copy-ratings-iframe-btn');
	if (generateBtn) {
		generateBtn.addEventListener('click', generateRatingsIframeCode);
	}
  if (copyBtn ) {
		copyBtn.addEventListener('click', copyToClipboard);
	}
}

/**
 * Generates an embeddable iframe code for ratings widget on a given URL.
 *
 * The user’s domain must be in the allowlist on the server for the iframe to function,
 * but validation is performed client-side for immediate feedback.
 */
export function generateRatingsIframeCode(): void {
	const urlInput = document.getElementById('rating-embed-url') as HTMLInputElement;
	const uiSelect = document.getElementById('rating-embed-ui') as HTMLSelectElement;
	const form = urlInput.closest('form') as HTMLFormElement;
	if (!form.checkValidity()) {
		urlInput.reportValidity();
		return;
	}

	const url = urlInput.value.trim();
	const ui = uiSelect.value || 'thumbs';

	const baseUrl = window.location.origin;
	const params = new URLSearchParams();
	params.set('url', url);
	params.set('ui', ui);

	const iframeUrl = `${baseUrl}/rating/v1?${params.toString()}`;

	// Calculate dimensions based on UI type
	let width = 120;
	let height = 50;
	if (ui === 'stars') {
		width = 200;
		height = 50;
	}

	const iframeCode = `<iframe src="${iframeUrl}" width="${width}" height="${height}" style="border:0;overflow:hidden;" loading="lazy" referrerpolicy="no-referrer"></iframe>`;

	(document.getElementById('ratings-iframe-code') as HTMLTextAreaElement).value = iframeCode;
	document.getElementById('ratings-embed-result')!.classList.remove('hidden');
}

/**
 * Copies the generated iframe code to the clipboard and provides visual feedback.
 * Temporarily changes the button text to "Copied!" for 2 seconds to provide user feedback,
 * then restores the original button text.
 */
function copyToClipboard(event: Event): void {
	const textarea = document.getElementById('ratings-iframe-code') as HTMLTextAreaElement;
	textarea.select();
	document.execCommand('copy');

	const btn = event.target as HTMLButtonElement;
	const originalText = btn.textContent;
	btn.textContent = 'Copied!';
	setTimeout(() => {
		btn.textContent = originalText;
	}, 2000);
}
