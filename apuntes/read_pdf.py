import sys
import pypdf

def main():
    pdf_path = r"c:\Users\jaime\OneDrive - Institución Educativa SEK\UNI\año2\semestre 2\REDES Y SISTEMAS WEB\RSW-PEC1_Grupo11-Jaime-Izan\apuntes\T3 - Tecnologías de servidor.pdf"
    try:
        reader = pypdf.PdfReader(pdf_path)
        for page in reader.pages:
            print(page.extract_text())
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    main()
