import fasttext.util
import csv


def get_words_list(path: str) -> [str]:
    words = list()
    with open(path, newline='\n') as csvfile:
        reader = csv.reader(csvfile, delimiter=',')
        for row in reader:
            words.append(', '.join(row))
    return words


fasttext.util.download_model('en', if_exists='ignore')
ft = fasttext.load_model('cc.en.300.bin')
print(ft.get_dimension())

path_to_csv = 'data.csv'
list_words = get_words_list(path_to_csv)
print(list_words)
