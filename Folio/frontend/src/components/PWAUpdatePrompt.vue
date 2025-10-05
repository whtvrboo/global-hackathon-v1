<template>
    <div v-if="showUpdatePrompt" class="fixed bottom-4 right-4 z-50">
        <div
            class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg p-4 max-w-sm">
            <div class="flex items-start">
                <div class="flex-shrink-0">
                    <svg class="w-5 h-5 text-blue-500" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd"
                            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                            clip-rule="evenodd"></path>
                    </svg>
                </div>
                <div class="ml-3 flex-1">
                    <h3 class="text-sm font-medium text-gray-900 dark:text-white">
                        Update Available
                    </h3>
                    <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                        A new version of Folio is available. Refresh to get the latest features.
                    </p>
                    <div class="mt-3 flex space-x-2">
                        <button @click="updateApp"
                            class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-2 px-3 rounded-md transition-colors">
                            Update Now
                        </button>
                        <button @click="dismissPrompt"
                            class="bg-gray-200 hover:bg-gray-300 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200 text-xs font-medium py-2 px-3 rounded-md transition-colors">
                            Later
                        </button>
                    </div>
                </div>
                <div class="ml-4 flex-shrink-0">
                    <button @click="dismissPrompt"
                        class="bg-white dark:bg-gray-800 rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">
                        <span class="sr-only">Close</span>
                        <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd"
                                d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                                clip-rule="evenodd"></path>
                        </svg>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref } from 'vue'

export default {
    name: 'PWAUpdatePrompt',
    setup() {
        const showUpdatePrompt = ref(false)
        let updateSW = null

        // This will be called from main.js when an update is available
        window.showPWAUpdatePrompt = (updateServiceWorker) => {
            updateSW = updateServiceWorker
            showUpdatePrompt.value = true
        }

        const updateApp = () => {
            if (updateSW) {
                updateSW(true) // Skip waiting and activate the new service worker
                showUpdatePrompt.value = false
            }
        }

        const dismissPrompt = () => {
            showUpdatePrompt.value = false
        }

        return {
            showUpdatePrompt,
            updateApp,
            dismissPrompt
        }
    }
}
</script>
