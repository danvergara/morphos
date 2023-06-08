# Python program to convert
# text file to pdf file


import cchardet
from fpdf import FPDF


class File:
    def __init__(self, file_name):
        self.file_name = file_name
        #! maybe this attributes should be handled
        #! inside the class methods
        # self.file_path = file_path
        # self.file_output_path = file_output_path
        # self.file_format = file_format

    def convert_from_doc_to_pdf(self):
        pdf = FPDF()
        pdf.add_page()
        pdf.set_font("Times", size=12)
        with open(self.file_name, "rb") as input_file:
            for x in input_file:
                detection = cchardet.detect(x)
                print(detection)
                detection_value = detection.get("encoding", "utf-8")
                if detection_value is None:
                    detection_value = "utf-8-sig"
                print(detection_value)
                pdf.cell(200, 10, txt=x.decode(detection_value), ln=1, align="C")

        pdf.output(f"{self.file_name}.pdf")


file_test = File("lorem_ipsum.docx")
file_test.convert_from_doc_to_pdf()
