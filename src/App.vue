<script setup>
import { ref, computed } from 'vue'

const originalPrompt = ref('')
const refinedPrompt = ref('')
const isLoading = ref(false)

const characterCount = computed(() => originalPrompt.value.length)
const hasPrompt = computed(() => originalPrompt.value.trim().length > 0)

const examplePrompts = [
  "Напиши код для сортировки массива",
  "Объясни квантовую физику простыми словами", 
  "Создай план тренировок на неделю",
  "Помоги написать резюме для программиста"
]

const refinePrompt = async () => {
  if (!hasPrompt.value) return
  
  isLoading.value = true
  
  try {
    const response = await fetch("http://localhost:8080/generate", {
      method: "POST",
      headers: { 
        "Content-Type": "application/json" 
      },
      body: JSON.stringify({ 
        input: originalPrompt.value 
      })
    })
    
    if (!response.ok) {
      throw new Error(`Ошибка сервера: ${response.status}`)
    }    
    const data = await response.json()
    
    refinedPrompt.value = data.prompt
  } catch (error) {
    console.error("Ошибка при обращении к серверу:", error)
    refinedPrompt.value = `❌ Ошибка подключения к серверу:

${error.message}

Проверьте:
• Запущен ли Go-сервер на порту 8080
• Доступен ли http://localhost:8080/generate
• Подключение к интернету

Попробуйте еще раз через несколько секунд.`
    
  } finally {
    isLoading.value = false
  }}
const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(refinedPrompt.value)
    alert('Скопировано в буфер обмена!')
  } catch (err) {
    console.error('Ошибка копирования:', err)
  }
}
const useExample = (example) => {
  originalPrompt.value = example
}
</script>
<template>
 <div id="app">
 <header>
 <img src="./assets/logoh1.png" alt="RefinedAI Logo" class="logo" />
<h1>RefinedAI</h1>
 <p>Улучшаем промты для искусственного интеллекта</p>
</header>
<main>
 <section class="examples-section">
 <h3><img src="./assets/lamp.png" alt="Лампочка" class="lamp-icon" /> Примеры промтов:</h3>
 <div class="examples-grid">
  <button 
 v-for="(example, index) in examplePrompts" 
 :key="index"
  @click="useExample(example)"
  class="example-btn"
 >
 {{ example }}
 </button>
</div>
</section>
 <section class="main-panels">
  <div class="input-panel">
    <div class="panel-header">
   <h3>Исходный промт</h3>
   <span class="char-counter">{{ characterCount }} символов</span>
 </div>
          
      <textarea
       v-model="originalPrompt"
  placeholder="Введите ваш промт здесь...
            
Например:
• Напиши функцию для...
• Объясни как работает...
• Создай план для..."
            rows="15"
            class="prompt-input"
 ></textarea>
          
 <div class="input-actions">
   <button 
        @click="refinePrompt" 
    :disabled="!hasPrompt || isLoading"
    class="refine-btn"
  >
 {{ isLoading ? 'Обрабатываем...' : 'Улучшить промт' }}
  </button>
</div>
</div>
  <div class="output-panel">
<div class="panel-header">
  <h3>Улучшенный промт</h3>
  <button 
 v-if="refinedPrompt" 
 @click="copyToClipboard"
 class="copy-btn"
     >
  Копировать
 </button>
   </div>
   <div class="prompt-output">
   <div v-if="isLoading" class="loading">
   <div class="spinner">⟳</div>
   <p>Анализируем ваш промт...</p>
   </div>   
    <pre v-else-if="refinedPrompt" class="refined-text">{{ refinedPrompt }}</pre>
    <div v-else class="placeholder">
    <p>Здесь появится ваш улучшенный промт</p>
  <small>Введите исходный промт и нажмите "Улучшить"</small>
  </div>
  </div>
 </div>
      </section></main>

  </div>
</template>

<style scoped>
.input-actions {
  display: flex;
  justify-content: left;
  margin-top: 1rem;
  padding: 0 1rem 1rem 1rem;
  position: relative;
  z-index: 3;
}

.refine-btn {
  padding: 0.6rem 1.2rem;
  border: 1px solid #8e8e93;
  border-radius: 6px;
  background: #2d2d2d;
  color: #ffffff;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.8rem;
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
}

.refine-btn:hover:not(:disabled) {
  background: #3a3a3a;
  border-color: #aeaeb2;
}

.refine-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.8rem 1.2rem;
  margin: 0 -5px -5px -5px;
  background: #2d2d2d;
  border-bottom: 1px solid #3a3a3a;
  position: relative;
  z-index: 3;
  border-top-left-radius: 12px;
  border-top-right-radius: 12px;
}



.panel-header h3 {
  margin: 0;
  color: #ffffff;
  font-size: 0.9rem;
  font-weight: 600;
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.char-counter {
  color: #8e8e93;
  font-size: 0.75rem;
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
  background: #1e1e1e;
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  border: 1px solid #3a3a3a;
}

.prompt-input {
  width: 100%;
  max-width: 100%;
  background: #1e1e1e;
  border: none;
  border-radius: 0;
  padding: 1rem;
  color: #ffffff;
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
  font-size: 0.9rem;
  line-height: 1.5;
  resize: none;
  height: 350px;
  overflow-y: auto;
  transition: all 0.3s ease;
  position: relative;
  z-index: 3;
  word-wrap: normal;
  word-break: normal;
  white-space: pre-wrap;
  overflow-wrap: break-word;
  overflow-x: hidden;
}

/* Адаптивная ширина для разных разрешений */
@media (min-width: 1920px) {
  .prompt-input {
    width: 120ch;
    max-width: calc(100% - 2rem);
  }
}

@media (min-width: 1440px) and (max-width: 1919px) {
  .prompt-input {
    width: 100ch;
    max-width: calc(100% - 2rem);
  }
}

@media (min-width: 1024px) and (max-width: 1439px) {
  .prompt-input {
    width: 80ch;
    max-width: calc(100% - 2rem);
  }
}

@media (min-width: 768px) and (max-width: 1023px) {
  .prompt-input {
    width: 60ch;
    max-width: calc(100% - 2rem);
  }
}

@media (max-width: 767px) {
  .prompt-input {
    width: 100%;
    max-width: calc(100% - 2rem);
    font-size: 0.8rem;
  }
}

.prompt-input:focus {
  outline: none;
  background: #1e1e1e;
}

.prompt-input::placeholder {
  color: #8e8e93;
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
}

.prompt-input::-webkit-scrollbar {
  width: 6px;
}

.prompt-input::-webkit-scrollbar-track {
  background: transparent;
}

.prompt-input::-webkit-scrollbar-thumb {
  background: rgba(142, 142, 147, 0.3);
  border-radius: 8px;
  border: 1px solid transparent;
  background-clip: content-box;
}

.prompt-input::-webkit-scrollbar-thumb:hover {
  background: rgba(142, 142, 147, 0.6);
}

.prompt-input::-webkit-scrollbar-thumb:active {
  background: rgba(142, 142, 147, 0.8);
}

.prompt-output::-webkit-scrollbar {
  width: 6px;
}

.prompt-output::-webkit-scrollbar-track {
  background: transparent;
}

.prompt-output::-webkit-scrollbar-thumb {
  background: rgba(142, 142, 147, 0.3);
  border-radius: 8px;
  border: 1px solid transparent;
  background-clip: content-box;
}

.prompt-output::-webkit-scrollbar-thumb:hover {
  background: rgba(142, 142, 147, 0.6);
}

.prompt-output::-webkit-scrollbar-thumb:active {
  background: rgba(142, 142, 147, 0.8);
}

.input-actions {
  display: flex;
  gap: 0.8rem;
  margin-top: 1rem;
  padding: 0 1rem 1rem 1rem;
  position: relative;
  z-index: 3;
}

.refine-btn, .clear-btn {
  padding: 0.6rem 1.2rem;
  border: 1px solid #8e8e93;
  border-radius: 6px;
  background: #2d2d2d;
  color: #ffffff;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.8rem;
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
}

.refine-btn:hover:not(:disabled) {
  background: #3a3a3a;
  border-color: #aeaeb2;
}

.clear-btn:hover:not(:disabled) {
  background: #3a3a3a;
  border-color: #aeaeb2;
}

.refine-btn:disabled, .clear-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.prompt-output {
  padding: 1rem;
  height: 350px;
  background: #1e1e1e;
  position: relative;
  z-index: 3;
  overflow-y: auto;
}

.placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: #8e8e93;
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
  text-align: center;
}

.placeholder-icon {
  font-size: 2rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.refined-text {
  color: #ffffff;
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
  font-size: 0.9rem;
  line-height: 1.5;
  white-space: pre-wrap;
  margin: 0;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: #8e8e93;
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
}

.spinner {
  font-size: 2rem;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.copy-btn {
  padding: 0.4rem 0.8rem;
  border: 1px solid #8e8e93;
  border-radius: 4px;
  background: #2d2d2d;
  color: #ffffff;
  cursor: pointer;
  font-size: 0.7rem;
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
  transition: all 0.2s ease;
}

.copy-btn:hover {
  background: #3a3a3a;
  border-color: #aeaeb2;
}
.main-panels {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin: 3rem 0;
  font-family: 'SF Mono', 'Monaco', 'Consolas', 'Liberation Mono', 'Courier New', monospace;
}

.input-panel, .output-panel {
  background: #1e1e1e;
  border: 1px solid #3a3a3a;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 
    0 8px 25px rgba(0, 0, 0, 0.3),
    0 4px 10px rgba(0, 0, 0, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.05);
  transition: all 0.3s ease;
  position: relative;
}

.input-panel:hover, .output-panel:hover {
  box-shadow: 
    0 12px 35px rgba(0, 0, 0, 0.4),
    0 6px 15px rgba(0, 0, 0, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.08);
  transform: translateY(-2px);
}



.input-panel::after, .output-panel::after {
  content: "";
  position: absolute;
  top: 10px;
  left: 15px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #ff5f57;
  box-shadow: 
    20px 0 0 #ffbd2e,
    40px 0 0 #28ca42;
  z-index: 2;
}
.examples-section {
  margin: 3rem 0;
  text-align: center;
}

.examples-section h3 {
  color: #8b5cf6;
  margin-bottom: 2rem;
  font-size: 1.4rem;
  text-shadow: 0 0 10px rgba(139, 92, 246, 0.5);
}

.examples-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 1.5rem;
  justify-content: center;
  align-items: center;
}

.example-btn {
  padding: 1.2rem 1.5rem;
  background: rgba(139, 92, 246, 0.1);
  color: #8b5cf6;
  border: 2px solid rgba(139, 92, 246, 0.3);
  border-radius: 50px;
  cursor: pointer;
  transition: all 0.4s ease;
  backdrop-filter: blur(10px);
  position: relative;
}

.example-btn:hover {
  transform: scale(1.05);
  background: rgba(139, 92, 246, 0.2);
  border-color: #8b5cf6;
  box-shadow: 0 0 20px rgba(139, 92, 246, 0.4);
}
header p {
  font-size: 1rem;
  font-weight: 500;
  color: #6b7280;
  text-align: center;
  margin: 1rem 0 2.5rem 0;
  text-transform: uppercase;
  letter-spacing: 2px;
  font-family: 'Arial', sans-serif;
}
h1 {
  font-size: 3.8rem;
  font-weight: 900;
  color: #8b5cf6;
  text-align: center;
  margin: 0;
}

.logo {
  height: 7.2rem;
  width: auto;
  display: block;
  margin: 0 auto 1rem auto;
}

  .lamp-icon {
    height: 2em;
    width: auto;
    vertical-align: -0.5em;
    margin-right: 0;
    display: inline;
    background: none;
    border: none;
    image-rendering: -webkit-optimize-contrast;
    image-rendering: crisp-edges;}
</style>
