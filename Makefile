create-venv:
	uv venv --python 3.10

run-venv:
	source .venv/bin/activate

install-libraries:
	uv pip install -r requirements.txt

run:
	python main.py
