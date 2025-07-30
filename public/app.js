// Filter produk favorit
// (script pencarian menu)
document.addEventListener('DOMContentLoaded', function() {
  const searchInput = document.getElementById('menuSearch');
  if (searchInput) {
    searchInput.addEventListener('input', function() {
      const keyword = this.value.toLowerCase();
      document.querySelectorAll('#about-menu .col-12.col-md-3').forEach(function(card) {
        const nama = card.getAttribute('data-nama').toLowerCase();
        if (nama.includes(keyword)) {
          card.style.display = '';
        } else {
          card.style.display = 'none';
        }
      });
    });
  }

  // Login admin & navbar logic
  const loginBtn = document.getElementById('adminLoginBtn');
  const loginModal = new bootstrap.Modal(document.getElementById('adminLoginModal'));
  const loginForm = document.getElementById('adminLoginForm');
  const loginError = document.getElementById('loginError');
  const gearDropdownMenu = document.getElementById('gearDropdownMenu');
  const navProduk = document.querySelector('a.nav-link[href="products"]')?.parentElement;
  const navKategori = document.querySelector('a.nav-link[href="categories"]')?.parentElement;

  function updateGearDropdown() {
    if (localStorage.getItem('isAdmin') === 'true') {
      if (!document.getElementById('adminLogoutBtn')) {
        const logoutLi = document.createElement('li');
        logoutLi.innerHTML = '<a class="dropdown-item text-danger" href="#" id="adminLogoutBtn">Logout Admin</a>';
        gearDropdownMenu.appendChild(logoutLi);
      }
      document.getElementById('adminLoginBtn').style.display = 'none';
    } else {
      const logoutBtn = document.getElementById('adminLogoutBtn');
      if (logoutBtn) logoutBtn.parentElement.remove();
      document.getElementById('adminLoginBtn').style.display = '';
    }
  }

  function updateNavbar() {
    if (localStorage.getItem('isAdmin') === 'true') {
      if (navProduk) navProduk.style.display = '';
      if (navKategori) navKategori.style.display = '';
    } else {
      if (navProduk) navProduk.style.display = 'none';
      if (navKategori) navKategori.style.display = 'none';
    }
  }

  async function checkAdminSession() {
    try {
      const res = await fetch('/admin/check', { credentials: 'same-origin' });
      if (res.ok) {
        localStorage.setItem('isAdmin', 'true');
      } else {
        localStorage.removeItem('isAdmin');
      }
    } catch {
      localStorage.removeItem('isAdmin');
    }
    updateNavbar();
    updateGearDropdown();
  }

  checkAdminSession();

  if (loginBtn) {
    loginBtn.addEventListener('click', function(e) {
      e.preventDefault();
      loginModal.show();
      loginError.style.display = 'none';
      loginForm.reset();
    });
  }

  if (loginForm) {
    loginForm.addEventListener('submit', async function(e) {
      e.preventDefault();
      const username = document.getElementById('adminUsername').value.trim();
      const password = document.getElementById('adminPassword').value.trim();
      try {
        const res = await fetch('/admin/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username, password }),
          credentials: 'same-origin'
        });
        if (res.ok) {
          await checkAdminSession();
          loginModal.hide();
          Swal.fire({
            icon: 'success',
            title: 'Login admin berhasil!',
            showConfirmButton: false,
            timer: 1500
          });
        } else {
          loginError.textContent = 'Username atau password salah!';
          loginError.style.display = 'block';
        }
      } catch (err) {
        loginError.textContent = 'Terjadi kesalahan server!';
        loginError.style.display = 'block';
      }
    });
  }

  if (gearDropdownMenu) {
    gearDropdownMenu.addEventListener('click', async function(e) {
      if (e.target && e.target.id === 'adminLogoutBtn') {
        e.preventDefault();
        await fetch('/admin/logout', { method: 'POST', credentials: 'same-origin' });
        localStorage.removeItem('isAdmin');
        updateNavbar();
        updateGearDropdown();
      }
    });
  }
});
// Script filterProducts (jika dipakai di halaman lain)
function filterProducts(category) {
  const items = document.querySelectorAll('.product-item');
  items.forEach(item => {
    if (category === 'all' || item.getAttribute('data-category') === category) {
      item.style.display = 'block';
    } else {
      item.style.display = 'none';
    }
  });
} 