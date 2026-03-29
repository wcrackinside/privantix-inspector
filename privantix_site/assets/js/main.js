document.addEventListener('DOMContentLoaded', () => {
  const y = document.getElementById('year');
  if (y) y.textContent = new Date().getFullYear();
});

function showToast(message){
  const toastEl = document.getElementById('liveToast');
  if (!toastEl) return;
  const body = toastEl.querySelector('.toast-body');
  if (body && message) body.textContent = message;
  const toast = new bootstrap.Toast(toastEl);
  toast.show();
}

function copyEmail(e){
  if (e) e.preventDefault();
  const email = 'hola@privantix.io';
  navigator.clipboard?.writeText(email);
  showToast('Email copiado: ' + email);
}

function waitlistSubmit(e){
  e.preventDefault();
  showToast('¡Listo! Te agregamos a la lista de espera (demo). En producción, conecta esto a tu backend.');
}
