import os
import time
import json
from datetime import datetime

from flask import Flask, render_template, request, jsonify
import requests

import regex_google

app = Flask(__name__)

ITINARY = {'London': {'lat': [51.5074], 'long': [-0.1278]},
    'Edinburgh': {'lat': [55.9533], 'long': [-3.1883]},
    'Manchester': {'lat': [53.4808], 'long': [-2.2426]},
    'Birmingham': {'lat': [52.4862], 'long': [-1.8904]},
    'Glasgow': {'lat': [55.8642], 'long': [-4.2518]},
    'Liverpool': {'lat': [53.4084], 'long': [-2.9916]},
    'Bristol': {'lat': [51.4545], 'long': [-2.5879]},
    'Cambridge': {'lat': [52.2053], 'long': [0.1218]},
    'Oxford': {'lat': [51.7548], 'long': [-1.2544]},
    'Leeds': {'lat': [53.8008], 'long': [-1.5491]},
    'York': {'lat': [53.959], 'long': [-1.0815]},
    'Brighton': {'lat': [50.8225], 'long': [-0.1372]},
    'Newcastle': {'lat': [54.9783], 'long': [-1.6174]},
    'Cardiff': {'lat': [51.4816], 'long': [-3.1791]},
    'Bath': {'lat': [51.3811], 'long': [-2.359]},
    'Belfast': {'lat': [54.5973], 'long': [-5.9301]},
    'Norwich': {'lat': [52.6309], 'long': [1.2974]},
    'Sheffield': {'lat': [53.3811], 'long': [-1.4701]},
    'Nottingham': {'lat': [52.9548], 'long': [-1.1581]},
    'Aberdeen': {'lat': [57.1497], 'long': [-2.0943]},
    'Portsmouth': {'lat': [50.8198], 'long': [-1.088]},
    'Southampton': {'lat': [50.9097], 'long': [-1.4044]},
    'Inverness': {'lat': [57.4778], 'long': [-4.2247]},
    'Stirling': {'lat': [56.1165], 'long': [-3.9369]},
    'Exeter': {'lat': [50.7184], 'long': [-3.5339]},
    'Chester': {'lat': [53.1915], 'long': [-2.8952]},
    'Dundee': {'lat': [56.462], 'long': [-2.9707]},
    'Canterbury': {'lat': [51.2802], 'long': [1.0789]},
    'Coventry': {'lat': [52.4066], 'long': [-1.5122]},
    'Swansea': {'lat': [51.6214], 'long': [-3.9436]},
    'Leicester': {'lat': [52.6369], 'long': [-1.1398]},
    'Luton': {'lat': [51.8787], 'long': [-0.42]},
    'Reading': {'lat': [51.4543], 'long': [-0.9781]},
    'Hull': {'lat': [53.7457], 'long': [-0.3367]},
    'Plymouth': {'lat': [50.3755], 'long': [-4.1427]},
    'Wolverhampton': {'lat': [52.5862], 'long': [-2.1269]},
    'Preston': {'lat': [53.7632], 'long': [-2.7031]},
    'Durham': {'lat': [54.7753], 'long': [-1.5849]},
    'Winchester': {'lat': [51.0626], 'long': [-1.3149]},
    'Salisbury': {'lat': [51.0688], 'long': [-1.7945]},
    'Lancaster': {'lat': [54.0484], 'long': [-2.7991]},
    'Ipswich': {'lat': [52.0567], 'long': [1.1482]},
    'Lincoln': {'lat': [53.2307], 'long': [-0.5406]},
    'Colchester': {'lat': [51.8892], 'long': [0.9042]},
    'Milton Keynes': {'lat': [52.0406], 'long': [-0.7594]},
    'Worcester': {'lat': [52.192], 'long': [-2.22]},
    'Peterborough': {'lat': [52.5695], 'long': [-0.2405]},
    'Gloucester': {'lat': [51.8642], 'long': [-2.2382]},
    'Bournemouth': {'lat': [50.7192], 'long': [-1.8808]},
    'St Albans': {'lat': [51.7527], 'long': [-0.3366]},
    'Llandudno': {'lat': [53.3247], 'long': [-3.8274]},
    'Scarborough': {'lat': [54.2817], 'long': [-0.4023]},
    'Whitby': {'lat': [54.4863], 'long': [-0.6145]},
    'Torquay': {'lat': [50.4619], 'long': [-3.5253]},
    'Margate': {'lat': [51.3813], 'long': [1.3862]},
    'Dover': {'lat': [51.1279], 'long': [1.3134]},
    'Sunderland': {'lat': [54.9069], 'long': [-1.3838]},
    'Blackpool': {'lat': [53.8175], 'long': [-3.0357]},
    'Windsor': {'lat': [51.4817], 'long': [-0.6117]},
    'Ely': {'lat': [52.3995], 'long': [0.2624]},
    'Rochester': {'lat': [51.388], 'long': [0.5089]},
    'Dover': {'lat': [51.1279], 'long': [1.3134]},
    'Folkestone': {'lat': [51.0814], 'long': [1.1698]},
    'Dartmoor': {'lat': [50.5668], 'long': [-3.9469]},
    'Snowdonia': {'lat': [53.0685], 'long': [-3.8257]},
    'Lake District': {'lat': [54.4609], 'long': [-3.0886]},
    'Cornwall': {'lat': [50.266], 'long': [-5.0527]},
    'Isle of Wight': {'lat': [50.6938], 'long': [-1.3047]},
    'Cotswolds': {'lat': [51.8248], 'long': [-1.6361]},
    'New Forest': {'lat': [50.8724], 'long': [-1.6009]},
    'Peak District': {'lat': [53.3437], 'long': [-1.8392]},
    'Stonehenge': {'lat': [51.1789], 'long': [-1.8262]},
    'Loch Ness': {'lat': [57.3229], 'long': [-4.4244]},
    'Ben Nevis': {'lat': [56.7969], 'long': [-5.0036]},
    'Skye': {'lat': [57.5359], 'long': [-6.2263]},
    'Isle of Mull': {'lat': [56.4686], 'long': [-5.9965]},
    'Edinburgh Castle': {'lat': [55.9486], 'long': [-3.1999]},
    'Tower of London': {'lat': [51.5081], 'long': [-0.0759]},
    'Buckingham Palace': {'lat': [51.5014], 'long': [-0.1419]},
    'British Museum': {'lat': [51.5194], 'long': [-0.127]},
    }

DATES_DICT = {'Today': 0,
              'Tomorrow': 1,
              '2 days from now': 2,
              '3 days from now': 3,
              '4 days from now': 4,
}

NUMBER_OF_DAYS = ["1", "2", "3", "4", "5", "6", "7"]

KELVIN_TO_CELCIUS = -273.15

global PLACES
PLACES=[]

@app.route("/")
def homepage():
    userPlaceChoice = request.args.get("Location", "London")

    if userPlaceChoice not in PLACES:
        print(f"\nUser Choice: |{userPlaceChoice}|") # This is to show the choice for debugging!
        PLACES.append(userPlaceChoice)

    Date = request.args.get("Date", "Today")
    Length_Of_Stay = request.args.get("Length_Of_Stay", "1")

    new_places = set(PLACES)
    location_infos = []
    print("Current list of places: ", PLACES) # This is to show the list of places in case their is a hiccup!

    if userPlaceChoice == '':
        userPlaceChoice = "British Museum"
    
    for place in new_places:
        existingPlace = place.title()
        if existingPlace in ITINARY.keys():
            Location_Info = get_info(Place=place.title(), coOrdinates=ITINARY[place], Length_Of_Stay=Length_Of_Stay, Date=Date)
            location_infos.append({"Location": place, "Info": Location_Info})

        elif 'google' in place:
            googlePlace = regex_google.getGoogleURL(place),
            lat = googlePlace[0][1]['lat'][0]
            lon = googlePlace[0][1]['long'][0]
            ITINARY[googlePlace[0][0]] = {'lat': [lat], 'long': [lon]}
            Location_Info = get_info(Place = googlePlace[0][0], coOrdinates= googlePlace[0][1], Length_Of_Stay=Length_Of_Stay, Date=Date)
            location_infos.append({"Location": googlePlace[0][0], "Info": Location_Info})

        else:
            # Check cache does the location exist there
            location_infos.append({"Location": place, "Info": [["Unable to identify location"]]})
            
    # Get a list of all available cities
    Cities = list(ITINARY.keys())

    # Debug printing
    # print(location_infos)

    # Render the template with multiple location information
    return render_template("locations.html",
                           Date=Date,
                           Length_Of_Stay=Length_Of_Stay,
                           DATES=DATES_DICT.keys(),
                           NUMBER_OF_DAYS=NUMBER_OF_DAYS,
                           Location_Infos=location_infos,
                           Cities=Cities)

@app.route('/clear_places', methods=['POST'])
def clear_places():
    global PLACES
    PLACES=['British Museum']
    print("Empty places", PLACES)
    return jsonify({"status": "success", "message": "PLACES array has been cleared"}), 200

def getKey():
    with open('key.txt', 'r') as f:
        key = f.read()
        key = key.strip()
    return key


def get_info(Place, coOrdinates, Length_Of_Stay, Date):
    print("Checking file exsits: ", Place)
    checkFileExists = checkFileExistance(Place, coOrdinates)

    with open(checkFileExists, 'r') as r:
        contents = r.read()

    weatherData = json.loads(contents)
    results = []
    dayGroup = []
    prev_day = datetime.fromtimestamp(weatherData['list'][0]['dt'])
    prev_day = prev_day.strftime('%d-%m-%Y') # yes this can be combined with the above but like this for readability
    count = weatherData['cnt']
    length_of_stay_number = 8
    too_big = False
    if (int(Length_Of_Stay) * 8) < 40:
        length_of_stay_number = int(Length_Of_Stay) * 8
        too_big = False
    else:
        length_of_stay_number = count
        too_big = True 
    x = 0 # just so X can't be unbound although i should always start as 0 as this is the start of the range, and logic is in place to make sure it only runs the first time and then can compound going forward!

    for i in range(0, length_of_stay_number):
        # Work out how many days to go through here
        if i == 0:
            x = 0
            offset_day = datetime.fromtimestamp(weatherData['list'][x]['dt'])
            looking_for_day = offset_day.strftime('%d')
            looking_for_day = DATES_DICT[Date] + int(looking_for_day)
            while looking_for_day != offset_day:
                offset_day = int(datetime.fromtimestamp(weatherData['list'][x]['dt']).strftime('%d'))
                # print(f'Offset day is: {offset_day}, Looking for day is: {looking_for_day}')
                if Date != "Today":
                    x += 1
        dayOffset = i + x
        print(f'Day offset is: {dayOffset} i is: {i} x is {x}')

        if dayOffset < len(weatherData['list']):
        # Getting all the details
            if int(str(datetime.fromtimestamp(weatherData['list'][dayOffset]['dt']).strftime('%H'))) < 5:
                continue
            weatherTime = str(datetime.fromtimestamp(weatherData['list'][dayOffset]['dt']).strftime('%H:%M'))
            day = datetime.fromtimestamp(weatherData['list'][dayOffset]['dt'])
            day = day.strftime('%d-%m-%Y')
            feelsLike = str(int(round(
                weatherData['list'][dayOffset]['main']['feels_like'], 0) + KELVIN_TO_CELCIUS)) + " ᵒC"
            weatherType = weatherData['list'][dayOffset]['weather'][0]['main']
            weatherDetail = weatherData['list'][dayOffset]['weather'][0]['description']
            windSpeed = round((weatherData['list'][dayOffset]['wind']['speed']) * 3.6, 2)

            if weatherType == "Rain":
                rainVolumeLast3Hours = weatherData['list'][dayOffset]['rain']['3h']
            else:
                rainVolumeLast3Hours = 0

            if day == prev_day:
                dayGroup.append({'weatherTime':weatherTime, 'feelsLike':feelsLike, 'weatherType':weatherType, 'weatherDetail':weatherDetail, 'rainVolume':rainVolumeLast3Hours, 'windSpeed': windSpeed})
            else:
                print(dayGroup)
                if dayGroup != []:
                    results.append({'timeDate':prev_day, 'listOfInfo':dayGroup})
                dayGroup = []
                prev_day = day

            print(day, {'weatherTime':weatherTime, 'feelsLike':feelsLike, 'weatherType':weatherType, 'weatherDetail':weatherDetail, 'rainVolume':rainVolumeLast3Hours, 'windSpeed': windSpeed})

    if dayGroup != []:
        results.append({'timeDate':prev_day, 'listOfInfo':dayGroup})

    # capture last day
    if too_big:
        results.append({'timeDate':prev_day, 'listOfInfo':dayGroup})
        results.append({'timeDate':"A future date - I am unable to predict beyond this point", 'listOfInfo':[]})

    return results


def checkFileExistance(Place, coOrdinates):
    checkFileExists = './caches/' + Place + '_cache.json'

    doesExist = os.path.exists(checkFileExists)
    
    if not doesExist:
        getNewLocation(coOrdinates, filePath=checkFileExists)
    else:
        max_age_hours = 3
        current_time = time.time()
        file_mod_time = os.path.getmtime(checkFileExists)
    
        # Calculate the file's age in seconds
        file_age_seconds = current_time - file_mod_time
    
        # Convert max_age_hours to seconds
        max_age_seconds = max_age_hours * 3600
    
        # Check if the file is older than the specified max age
        if file_age_seconds > max_age_seconds:
            os.remove(checkFileExists)
            getNewLocation(coOrdinates, checkFileExists)
    
    return checkFileExists

def getNewLocation(coOrdinates, filePath, key=getKey()):
    lat = coOrdinates['lat'][0]
    lon = coOrdinates['long'][0]
    website = f"https://api.openweathermap.org/data/2.5/forecast?lat={lat}&lon={lon}&appid={key}"
    response = requests.get(website)

    if response.status_code == 200:
        with open(filePath, 'w') as w:
            w.write(response.text)
    else:
        print("Status response as it's not ok!", response)

    time.sleep(1)

if __name__ == '__main__':
    # Replace '192.168.x.x' with your actual IP address
    app.run(host='192.168.1.75', port=8000, debug=True)
    # app.run(host='192.168.0.223', port=5000, debug=True)
