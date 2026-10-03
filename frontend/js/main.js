const API_URL = 'http://localhost:8080/api/v1';

// Элементы
const videoInput = document.getElementById('videoInput');
const uploadBtn = document.getElementById('uploadBtn');
const blurBtn = document.getElementById('blurBtn');
const downloadBtn = document.getElementById('downloadBtn');
const originalContent = document.getElementById('originalContent');
const blurredContent = document.getElementById('blurredContent');
const originalVideo = document.getElementById('originalVideo');
const blurredVideo = document.getElementById('blurredVideo');

// Состояние
let selectedFile = null;
let uploadedVideo = null;
let blurredVideoData = null;

// ============ 1. Upload ============
uploadBtn.addEventListener('click', () => videoInput.click());

videoInput.addEventListener('change', async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    selectedFile = file;

    // Показать локально
    originalVideo.src = URL.createObjectURL(file);
    originalVideo.hidden = false;
    originalContent.hidden = true;

    // Сбросить blur, если было предыдущее видео
    blurredVideo.hidden = true;
    blurredContent.hidden = false;
    blurBtn.disabled = true;
    blurBtn.textContent = 'Blur';
    downloadBtn.disabled = true;

    try {
        uploadBtn.disabled = true;
        uploadBtn.textContent = 'Loading...';

        const formData = new FormData();
        formData.append('video', file);

        const resp = await fetch(`${API_URL}/videos`, {
            method: 'POST',
            body: formData,
        });

        if (!resp.ok) throw new Error(await resp.text());

        uploadedVideo = await resp.json();

        blurBtn.disabled = false;
        uploadBtn.textContent = 'Uploaded';

    } catch (err) {
        alert('Upload error: ' + err.message);
        uploadBtn.textContent = 'Upload';
        uploadBtn.disabled = false;
    }
});

// ============ 2. Blur ============
blurBtn.addEventListener('click', async () => {
    if (!uploadedVideo) return;

    try {
        blurBtn.disabled = true;
        blurBtn.textContent = 'Processing...';

        const resp = await fetch(`${API_URL}/videos`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                id: uploadedVideo.id,
                version: uploadedVideo.version,
            }),
        });

        if (!resp.ok) throw new Error(await resp.text());

        blurredVideoData = await resp.json();

        // Показать заблюренное
        const fileURL = `${API_URL}/videos/${blurredVideoData.id}/file`;
        blurredVideo.src = fileURL;
        blurredVideo.hidden = false;
        blurredContent.hidden = true;

        downloadBtn.disabled = false;
        blurBtn.textContent = 'Done';

    } catch (err) {
        alert('Blur error: ' + err.message);
        blurBtn.textContent = 'Blur';
        blurBtn.disabled = false;
    }
});

// ============ 3. Download ============
downloadBtn.addEventListener('click', () => {
    if (!blurredVideoData) return;

    const fileURL = `${API_URL}/videos/${blurredVideoData.id}/file`;

    const a = document.createElement('a');
    a.href = fileURL;
    a.download = blurredVideoData.file_name;
    a.click();
});