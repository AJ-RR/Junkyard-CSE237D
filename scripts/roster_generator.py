import csv
import random

first_names = [
    "Alex", "Brandon", "Cameron", "Dylan", "Evan", "Felix", "Grace", "Hannah", "Isabel", "Jack",
    "Kara", "Liam", "Mia", "Nora", "Owen", "Paige", "Quinn", "Ryan", "Sophie", "Tyler"]
last_names = [
    "Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Martinez", "Hernandez",
    "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee"
]

def generate_email(first, last, idx):
    return f"{first[0].lower()}{last.lower()}@example.edu"

with open('students.csv', 'w', newline='') as csvfile:
    writer = csv.writer(csvfile)
    writer.writerow(['first_name', 'last_name', 'email'])
    for i in range(1, 101):
        first = random.choice(first_names)
        last = random.choice(last_names)
        email = generate_email(first, last, i)
        writer.writerow([first, last, email])