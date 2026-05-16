'use strict';

document.addEventListener('DOMContentLoaded', function () {
	const listener = new Listener();

	listener.decode = function () {
		const a = document.getElementById('link-interaction');

		if (!a) {
			return;
		}

		const encoded = a.getAttribute('href');
		const decoded = encoded.replace('-', '@').replaceAll('-', '.');

		a.setAttribute('href', 'mailto:' + decoded);
		a.textContent = decoded;
	};

	listener.on();
});

// Listener

function Listener() {}

Listener.prototype.decode = null;

Listener.prototype.on = function () {
	this.listener = this.__onInteraction.bind(this);

	document.addEventListener('mouseenter', this.listener, true);
	document.addEventListener('focus', this.listener, true);
	document.addEventListener('touchstart', this.listener, true);
};

Listener.prototype.off = function () {
	document.removeEventListener('mouseenter', this.listener, true);
	document.removeEventListener('focus', this.listener, true);
	document.removeEventListener('touchstart', this.listener, true);

	delete this.listener;
};

Listener.prototype.__onInteraction = function () {
	this.off();
	this.decode();
};
