from pdfminer.high_level import extract_text
import glob
x = glob.glob("pdf/*.pdf")

#ii=0
animals1 = []
animals2 = []

for i in x:
    #print(i)
    #print(i.replace("pdf\\", "",1))
    xi = i.replace("pdf\\", "pdf/",1)
   # xii = xi + ".txt"
    animals1.append(xi)
    #print(xi)
    #print("BREAk \n")

for i in x:
    #print(i)
    #print(i.replace("pdf\\", "",1))
    xi = i.replace("pdf\\", "",1)
    #print(xi)
    #print("BREAk \n")
    xi = xi.replace(".pdf", "",1)
    xii = xi + ".txt"
    animals2.append(xii)
    #print(xii)
    #print("BREAk \n")

#print("ANIMASL \n")
    #xi = xi.replace(".pdf", "",1)
#print(animals[0])

# text = extract_text(animals1[0])
# #print(extract_text(x[0]))
# print(text)



# with open('readme.txt', 'w') as f:
#     f.write('readme')

#text = extract_text('pdf/Options for Persistence of Cyberweapons.pdf')
print("ERGEBNIS")

for z,j in zip(animals1, animals2):
#     print(z)
#     print("BREAk \n")
    print(z)
    print("BREAk \n")
    # while True:
    #     try:
    #         with open(j, 'w', encoding="utf-8") as f:
    #             text = extract_text(z)
    #             f.write(text)
    #         break
    #     except (RuntimeError, TypeError, NameError):
    #         print("Oops!  That was no valid number.  Try again...")

    try:
        with open(j, 'w', encoding="utf-8") as f:
            text = extract_text(z)
            f.write(text)
    except OSError as err:
        print("OS error: {0}".format(err))
    except ValueError:
        print("Could not convert data to an integer.")
    except BaseException as err:
        print(f"Unexpected {err=}, {type(err)=}")
    #raise
    # while True:
    #     try:
    #         with open(j, 'w', encoding="utf-8") as f:
    #             text = extract_text(z)
    #             f.write(text)
    #         break
    #     except PDFSyntaxError:
    #         print("Oops!  That was no valid pdf.  Try again...")

    
        #encoding wirkt nur auf den txt
#print(animals)


# import pdfminer
# print(pdfminer.__version__)  