<template>
  <div class="auth-container">
    <div class="auth-card">
      <h2>Create Account</h2>

      <!-- Name -->
      <div class="form-group">
        <input
            v-model="form.name"
            type="text"
            placeholder="Full Name (2-50 characters)"
            class="input"
            @blur="validateField('name')"
        />
        <p v-if="errors.name" class="error-text">{{ errors.name }}</p>
      </div>

      <!-- Email -->
      <div class="form-group">
        <input
            v-model="form.email"
            type="email"
            placeholder="Email (example@domain.com)"
            class="input"
            @blur="validateField('email')"
        />
        <p v-if="errors.email" class="error-text">{{ errors.email }}</p>
      </div>

      <!-- Phone Number -->
      <div class="form-group">
        <input
            v-model="form.phone_number"
            type="tel"
            placeholder="Phone Number (10 digits)"
            class="input"
            @blur="validateField('phone_number')"
        />
        <p v-if="errors.phone_number" class="error-text">{{ errors.phone_number }}</p>
      </div>

      <!-- Address fields -->
      <div class="address-group">
        <div class="form-group">
          <select v-model="form.city" @change="onCityChange" class="input">
            <option value="" disabled>Select City</option>
            <option v-for="city in cities" :key="city.id" :value="city.name">{{ city.name }}</option>
          </select>
          <p v-if="errors.city" class="error-text">{{ errors.city }}</p>
        </div>

        <div class="form-group">
          <select v-model="form.district" @change="onDistrictChange" class="input" :disabled="!form.city">
            <option value="" disabled>Select District</option>
            <option v-for="district in filteredDistricts" :key="district.id" :value="district.name">{{ district.name }}</option>
          </select>
          <p v-if="errors.district" class="error-text">{{ errors.district }}</p>
        </div>

        <div class="form-group">
          <select v-model="form.ward" class="input" :disabled="!form.district">
            <option value="" disabled>Select Ward/Commune</option>
            <option v-for="ward in filteredWards" :key="ward.id" :value="ward.name">{{ ward.name }}</option>
          </select>
          <p v-if="errors.ward" class="error-text">{{ errors.ward }}</p>
        </div>

        <div class="form-group">
          <input
              v-model="form.street"
              type="text"
              placeholder="Street/Address Detail"
              class="input"
              @blur="validateField('street')"
          />
        <p v-if="errors.street" class="error-text">{{ errors.street }}</p>
        </div>
      </div>


      <!-- Password -->
      <div class="form-group">
        <input
            v-model="form.password"
            type="password"
            placeholder="Password (8-64 characters)"
            class="input"
            @blur="validateField('password')"
        />
        <p v-if="errors.password" class="error-text">{{ errors.password }}</p>
      </div>

      <!-- Role -->
      <div class="form-group">
        <select v-model="form.role" class="input" @blur="validateField('role')">
          <option value="">Select Role</option>
          <option value="customer">Customer</option>
          <option value="shipper">Shipper</option>
          <option value="restaurant_admin">Restaurant Admin</option>
        </select>
        <p v-if="errors.role" class="error-text">{{ errors.role }}</p>
      </div>

      <!-- Submit Button -->
      <button @click="handleRegister" class="btn btn-primary" :disabled="loading">
        {{ loading ? 'Creating Account...' : 'Register' }}
      </button>

      <!-- Footer -->
      <div class="auth-footer">
        <router-link to="/login">Already have an account? Login</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import {computed, ref} from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  name: '',
  email: '',
  phone_number: '',
  city: '',
  district: '',
  ward: '',
  street: '',
  password: '',
  role: ''
})

const errors = ref({
  name: '',
  email: '',
  phone_number: '',
  city: '',
  district: '',
  ward: '',
  street: '',
  password: '',
  role: ''
})

const cities = ref([
  {id: 1, name: 'Hà Nội'},
])

const districts = ref([
  { id: 1, cityId: 1, name: 'Ba Đình' },
  { id: 2, cityId: 1, name: 'Cầu Giấy' },
  { id: 3, cityId: 1, name: 'Hoàn Kiếm' },
  { id: 4, cityId: 1, name: 'Hai Bà Trưng' },
  { id: 5, cityId: 1, name: 'Hoàng Mai' },
  { id: 6, cityId: 1, name: 'Đống Đa' },
  { id: 7, cityId: 1, name: 'Tây Hồ' },
  { id: 8, cityId: 1, name: 'Thanh Xuân' },
  { id: 9, cityId: 1, name: 'Bắc Từ Liêm' },
  { id: 10, cityId: 1, name: 'Hà Đông' },
  { id: 11, cityId: 1, name: 'Long Biên' },
  { id: 12, cityId: 1, name: 'Nam Từ Liêm' },
])

const wards = ref([
  // Ba Đình
  { id: 1, districtId: 1, name: 'Cống Vị' },
  { id: 2, districtId: 1, name: 'Điện Biên' },
  { id: 3, districtId: 1, name: 'Đội Cấn' },
  { id: 4, districtId: 1, name: 'Giảng Võ' },
  { id: 5, districtId: 1, name: 'Kim Mã' },
  { id: 6, districtId: 1, name: 'Liễu Giai' },
  { id: 7, districtId: 1, name: 'Ngọc Hà' },
  { id: 8, districtId: 1, name: 'Ngọc Khánh' },
  { id: 9, districtId: 1, name: 'Nguyễn Trung Trực' },
  { id: 10, districtId: 1, name: 'Phúc Xá' },
  { id: 11, districtId: 1, name: 'Quán Thánh' },
  { id: 12, districtId: 1, name: 'Thành Công' },
  { id: 13, districtId: 1, name: 'Trúc Bạch' },
  { id: 14, districtId: 1, name: 'Vĩnh Phúc'},

  // Cầu Giấy
  { id: 15, districtId: 2, name: 'Dịch Vọng Hậu' },
  { id: 16, districtId: 2, name: 'Mai Dịch' },
  { id: 17, districtId: 2, name: 'Nghĩa Đô' },
  { id: 18, districtId: 2, name: 'Nghĩa Tân' },
  { id: 19, districtId: 2, name: 'Quan Hoa' },
  { id: 20, districtId: 2, name: 'Yên Hòa' },
  { id: 21, districtId: 2, name: 'Trung Hòa' },
  { id: 22, districtId: 2, name: 'Dịch Vọng' },

  // Hoàn Kiếm
  { id: 23, districtId: 3, name: 'Chương Dương' },
  { id: 24, districtId: 3, name: 'Cửa Đông' },
  { id: 25, districtId: 3, name: 'Cửa Nam' },
  { id: 26, districtId: 3, name: 'Đồng Xuân' },
  { id: 27, districtId: 3, name: 'Hàng Bạc' },
  { id: 28, districtId: 3, name: 'Hàng Bài' },
  { id: 29, districtId: 3, name: 'Hàng Bồ' },
  { id: 30, districtId: 3, name: 'Hàng Bông' },
  { id: 31, districtId: 3, name: 'Hàng Buồm' },
  { id: 32, districtId: 3, name: 'Hàng Đào' },
  { id: 33, districtId: 3, name: 'Hàng Gai' },
  { id: 34, districtId: 3, name: 'Hàng Mã' },
  { id: 35, districtId: 3, name: 'Hàng Trống' },
  { id: 36, districtId: 3, name: 'Lý Thái Tổ' },
  { id: 37, districtId: 3, name: 'Phan Chu Trinh' },
  { id: 38, districtId: 3, name: 'Phúc Tân' },
  { id: 39, districtId: 3, name: 'Trần Hưng Đạo' },
  { id: 40, districtId: 3, name: 'Tràng Tiền' },

  // Hai Bà Trưng
  { id: 41, districtId: 4, name: 'Quỳnh Lôi' },
  { id: 42, districtId: 4, name: 'Bách Khoa' },
  { id: 43, districtId: 4, name: 'Bạch Đằng' },
  { id: 44, districtId: 4, name: 'Cầu Dền' },
  { id: 45, districtId: 4, name: 'Đống Mác' },
  { id: 46, districtId: 4, name: 'Đồng Nhân' },
  { id: 47, districtId: 4, name: 'Đồng Tâm' },
  { id: 48, districtId: 4, name: 'Lê Đại Hành' },
  { id: 49, districtId: 4, name: 'Minh Khai' },
  { id: 50, districtId: 4, name: 'Nguyễn Du' },
  { id: 51, districtId: 4, name: 'Phạm Đình Hổ' },
  { id: 52, districtId: 4, name: 'Phố Huế' },
  { id: 53, districtId: 4, name: 'Quỳnh Mai' },
  { id: 54, districtId: 4, name: 'Thanh Lương' },
  { id: 55, districtId: 4, name: 'Thanh Nhàn' },
  { id: 56, districtId: 4, name: 'Trương Định' },
  { id: 57, districtId: 4, name: 'Vĩnh Tuy' },


  // Hoàng Mai
  { id: 58, districtId: 5, name: 'Đại Kim' },
  { id: 59, districtId: 5, name: 'Định Công' },
  { id: 60, districtId: 5, name: 'Giáp Bát' },
  { id: 61, districtId: 5, name: 'Hoàng Liệt' },
  { id: 62, districtId: 5, name: 'Hoàng Văn Thụ' },
  { id: 63, districtId: 5, name: 'Lĩnh Nam' },
  { id: 64, districtId: 5, name: 'Mai Động' },
  { id: 65, districtId: 5, name: 'Tân Mai' },
  { id: 66, districtId: 5, name: 'Thanh Trì' },
  { id: 67, districtId: 5, name: 'Thịnh Liệt' },
  { id: 68, districtId: 5, name: 'Trần Phú' },
  { id: 69, districtId: 5, name: 'Tương Mai' },
  { id: 70, districtId: 5, name: 'Vĩnh Hưng' },
  { id: 71, districtId: 5, name: 'Yên Sở' },

  // Đống Đa
  { id: 72, districtId: 6, name: 'Cát Linh' },
  { id: 73, districtId: 6, name: 'Hàng Bột' },
  { id: 74, districtId: 6, name: 'Khâm Thiên' },
  { id: 75, districtId: 6, name: 'Khương Thượng' },
  { id: 76, districtId: 6, name: 'Kim Liên' },
  { id: 77, districtId: 6, name: 'Láng Hạ' },
  { id: 78, districtId: 6, name: 'Láng Thượng' },
  { id: 79, districtId: 6, name: 'Nam Đồng' },
  { id: 80, districtId: 6, name: 'Ngã Tư Sở' },
  { id: 81, districtId: 6, name: 'Ô Chợ Dừa' },
  { id: 82, districtId: 6, name: 'Phương Liên' },
  { id: 83, districtId: 6, name: 'Phương Mai' },
  { id: 84, districtId: 6, name: 'Quang Trung' },
  { id: 85, districtId: 6, name: 'Quốc Tử Giám' },
  { id: 86, districtId: 6, name: 'Thịnh Quang' },
  { id: 87, districtId: 6, name: 'Thổ Quan' },
  { id: 88, districtId: 6, name: 'Trung Liệt' },
  { id: 89, districtId: 6, name: 'Trung Phụng' },
  { id: 90, districtId: 6, name: 'Trung Tự' },
  { id: 91, districtId: 6, name: 'Văn Chương' },
  { id: 92, districtId: 6, name: 'Văn Miếu' },


  // Tây Hồ
  { id: 93, districtId: 7, name: 'Bưởi' },
  { id: 94, districtId: 7, name: 'Nhật Tân' },
  { id: 95, districtId: 7, name: 'Phú Thượng' },
  { id: 96, districtId: 7, name: 'Quảng An' },
  { id: 97, districtId: 7, name: 'Thụy Khuê' },
  { id: 98, districtId: 7, name: 'Tứ Liên' },
  { id: 99, districtId: 7, name: 'Xuân La' },
  { id: 100, districtId: 7, name: 'Yên Phụ' },

  // Thanh Xuân
  { id: 101, districtId: 8, name: 'Hạ Đình' },
  { id: 102, districtId: 8, name: 'Khương Đình' },
  { id: 103, districtId: 8, name: 'Khương Mai' },
  { id: 104, districtId: 8, name: 'Khương Trung' },
  { id: 105, districtId: 8, name: 'Kim Giang' },
  { id: 106, districtId: 8, name: 'Nhân Chính' },
  { id: 107, districtId: 8, name: 'Phương Liệt' },
  { id: 108, districtId: 8, name: 'Thanh Xuân Bắc' },
  { id: 109, districtId: 8, name: 'Thanh Xuân Nam' },
  { id: 110, districtId: 8, name: 'Thanh Xuân Trung' },
  { id: 111, districtId: 8, name: 'Thượng Đình' },


  // Bắc Từ Liêm
  { id: 112, districtId: 9, name: 'Cổ Nhuế 1' },
  { id: 113, districtId: 9, name: 'Cổ Nhuế 2' },
  { id: 114, districtId: 9, name: 'Đông Ngạc' },
  { id: 115, districtId: 9, name: 'Đức Thắng' },
  { id: 116, districtId: 9, name: 'Liên Mạc' },
  { id: 117, districtId: 9, name: 'Minh Khai' },
  { id: 118, districtId: 9, name: 'Phú Diễn' },
  { id: 119, districtId: 9, name: 'Phúc Diễn' },
  { id: 120, districtId: 9, name: 'Tây Tựu' },
  { id: 121, districtId: 9, name: 'Thượng Cát' },
  { id: 122, districtId: 9, name: 'Thụy Phương' },
  { id: 123, districtId: 9, name: 'Xuân Đỉnh' },
  { id: 124, districtId: 9, name: 'Xuân Tảo' },


  // Hà Đông
  { id: 125, districtId: 10, name: 'Biên Giang' },
  { id: 126, districtId: 10, name: 'Đồng Mai' },
  { id: 127, districtId: 10, name: 'Dương Nội' },
  { id: 128, districtId: 10, name: 'Hà Cầu' },
  { id: 129, districtId: 10, name: 'Kiến Hưng' },
  { id: 130, districtId: 10, name: 'La Khê' },
  { id: 131, districtId: 10, name: 'Mộ Lao' },
  { id: 132, districtId: 10, name: 'Nguyễn Trãi' },
  { id: 133, districtId: 10, name: 'Phú La' },
  { id: 134, districtId: 10, name: 'Phú Lãm' },
  { id: 135, districtId: 10, name: 'Phú Lương' },
  { id: 136, districtId: 10, name: 'Phúc La' },
  { id: 137, districtId: 10, name: 'Quang Trung' },
  { id: 138, districtId: 10, name: 'Vạn Phúc' },
  { id: 139, districtId: 10, name: 'Văn Quán' },
  { id: 140, districtId: 10, name: 'Yên Nghĩa' },
  { id: 141, districtId: 10, name: 'Yết Kiêu' },


  // Long Biên
  { id: 142, districtId: 11, name: 'Bồ Đề' },
  { id: 143, districtId: 11, name: 'Cự Khối' },
  { id: 144, districtId: 11, name: 'Đức Giang' },
  { id: 145, districtId: 11, name: 'Gia Thụy' },
  { id: 146, districtId: 11, name: 'Giang Biên' },
  { id: 147, districtId: 11, name: 'Long Biên' },
  { id: 148, districtId: 11, name: 'Ngọc Lâm' },
  { id: 149, districtId: 11, name: 'Ngọc Thụy' },
  { id: 150, districtId: 11, name: 'Phúc Đồng' },
  { id: 151, districtId: 11, name: 'Phúc Lợi' },
  { id: 152, districtId: 11, name: 'Sài Đồng' },
  { id: 153, districtId: 11, name: 'Thạch Bàn' },
  { id: 154, districtId: 11, name: 'Thượng Thanh' },
  { id: 155, districtId: 11, name: 'Việt Hưng' },

  // Nam Từ Liêm
  { id: 156, districtId: 12, name: 'Cầu Diễn' },
  { id: 157, districtId: 12, name: 'Đại Mỗ' },
  { id: 158, districtId: 12, name: 'Mễ Trì' },
  { id: 159, districtId: 12, name: 'Mỹ Đình 1' },
  { id: 160, districtId: 12, name: 'Mỹ Đình 2' },
  { id: 161, districtId: 12, name: 'Phú Đô' },
  { id: 162, districtId: 12, name: 'Phương Canh' },
  { id: 163, districtId: 12, name: 'Tây Mỗ' },
  { id: 164, districtId: 12, name: 'Trung Văn' },
  { id: 165, districtId: 12, name: 'Xuân Phương' },

]);

const filteredDistricts = computed(() => {
  return districts.value.filter(d => d.cityId === cities.value.find(c => c.name === form.value.city)?.id)
})

const filteredWards = computed(() => {
  return wards.value.filter(w => w.districtId === districts.value.find(d => d.name === form.value.district)?.id)
})

const loading = ref(false)

//Reset downstream fields
function onCityChange() {
  form.value.district = ''
  form.value.ward = ''
}

function onDistrictChange() {
  form.value.ward = ''
}

// Validate từng field riêng
const validateField = (field) => {
  switch(field) {
    case 'name':
      if (!form.value.name) errors.value.name = 'Name is required'
      else if (form.value.name.length < 2 || form.value.name.length > 50) errors.value.name = 'Name must be 2-50 characters'
      else errors.value.name = ''
      break
    case 'email':
      const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!form.value.email) errors.value.email = 'Email is required'
      else if (!emailPattern.test(form.value.email)) errors.value.email = 'Invalid email format'
      else errors.value.email = ''
      break
    case 'phone_number':
      const phonePattern = /^\d{10}$/
      if (!form.value.phone_number) errors.value.phone_number = 'Phone number is required'
      else if (!phonePattern.test(form.value.phone_number)) errors.value.phone_number = 'Phone number must be 10 digits'
      else errors.value.phone_number = ''
      break
    case 'city':
      errors.value.city = form.value.city ? '' : 'City is required'
      break
    case 'district':
      errors.value.district = form.value.district ? '' : 'District is required'
      break
    case 'ward':
      errors.value.ward = form.value.ward ? '' : 'Ward is required'
      break
    case 'street':
      errors.value.street = form.value.street ? '' : 'Street is required'
      break
    case 'password':
      if (!form.value.password) errors.value.password = 'Password is required'
      else if (form.value.password.length < 8 || form.value.password.length > 64) errors.value.password = 'Password must be 8-64 characters'
      else errors.value.password = ''
      break
    case 'role':
      errors.value.role = form.value.role ? '' : 'Role is required'
      break
  }
}

// Validate toàn bộ form
const validateForm = () => {
  validateField('name')
  validateField('email')
  validateField('phone_number')
  validateField('city')
  validateField('district')
  validateField('ward')
  validateField('street')
  validateField('password')
  validateField('role')

  return Object.values(errors.value).every(e => !e)
}

const handleRegister = async () => {
  if (!validateForm()) return

  loading.value = true
  try {
    const fullAddress = `${form.value.street}, ${form.value.ward}, ${form.value.district}, ${form.value.city}`
    const payload = { ...form.value, address: fullAddress }
    delete payload.city
    delete payload.district
    delete payload.ward
    delete payload.street

    await authStore.register(payload)
    alert('Registration successful! Please login.')
    router.push('/login')
  } catch (error) {
    alert('Registration failed: ' + (error.response?.data?.message || error.message))
  } finally {
    loading.value = false
  }
}
</script>


<style scoped>
.auth-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #fb923c 0%, #ef4444 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.auth-card {
  background: white;
  border-radius: 1rem;
  padding: 2rem;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}

.auth-card h2 {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 1.5rem;
  text-align: center;
}

.form-group {
  margin-bottom: 1rem;
}

.input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.5rem;
  font-size: 1rem;
}

.input:focus {
  outline: none;
  border-color: #ea580c;
  box-shadow: 0 0 0 3px rgba(234, 88, 12, 0.1);
}

.btn {
  width: 100%;
  padding: 0.75rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #ea580c;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #c2410c;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.auth-footer {
  text-align: center;
  margin-top: 1.5rem;
}

.auth-footer a {
  color: #ea580c;
  text-decoration: none;
}

.auth-footer a:hover {
  text-decoration: underline;
}

.address-group {
  display: grid;
  grid-template-columns: repeat(2, 1fr); /* 2 cột bằng nhau */
  gap: 1rem; /* khoảng cách giữa các input */
}

.address-group .form-group {
  margin-bottom: 0; /* đã có gap nên không cần margin-bottom nữa */
}

/* Màu đỏ cho lỗi */
.error-text {
  color: red;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}
</style>
