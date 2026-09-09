# LETS Foam Calculator

**Version v6**

LETS Foam Calculator is a lightweight web application for estimating Class B firefighting foam-solution requirements. It is intended for training, planning, and instructional use by LETS Training Solutions.

The application is implemented as a small Go HTTP server with the calculator interface embedded directly in `main.go`. It can run locally or be deployed as a web service, including the project's Render deployment.

## What the calculator does

The calculator estimates three quantities:

1. **Total Foam Solution** — the total finished foam solution required for the selected spill area, application rate, and design duration.
2. **Foam Concentrate Needed** — the amount of foam concentrate contained in that finished solution at the selected concentration percentage.
3. **Water Needed** — the remaining water portion of the finished foam solution.

The normal workflow is intentionally simple for students. The user selects the liquid type, incident type, foam concentration, and spill size. The calculator automatically loads standard application-rate and design-duration values for the selected hazard combination.

## Primary inputs

The calculator accepts the following primary inputs:

- **Liquid Type**
  - Hydrocarbon
  - Polar Solvent
- **Incident Type**
  - Spill Fires
  - In Depth Fires
- **Foam Concentration**
  - 3%
  - 6%
- **Spill Size**
  - Entered in square feet
  - Current interface range: 100 to 100,000 sq ft

## Dynamic fire-service defaults

Version v6 adds automatic preloading of the Application Rate and Design Duration based on the selected Liquid Type and Incident Type.

| Liquid Type | Incident Type | Default Application Rate | Default Design Duration |
| --- | --- | ---: | ---: |
| Hydrocarbon | Spill Fire | 0.10 GPM/sq ft | 15 minutes |
| Hydrocarbon | In-Depth Fire | 0.10 GPM/sq ft | 30 minutes |
| Polar Solvent | Spill Fire | 0.10 GPM/sq ft | 15 minutes |
| Polar Solvent | In-Depth Fire | 0.20 GPM/sq ft | 60 minutes |

Whenever the user changes either Liquid Type or Incident Type, the application reloads the appropriate pair of default values.

## Advanced Options

Version v6 adds a collapsible **Advanced Options** card. It is closed by default so the normal student workflow remains uncluttered.

Opening Advanced Options exposes:

- **Application Rate (GPM/sq ft)**
- **Design Duration (minutes)**

These fields contain the currently preloaded values for the selected hazard combination.

### Manual overrides

Advanced users can type over either preloaded value.

The calculator does not substitute the default values again when the Calculate button is pressed. It reads the values currently present in the fields. This means a valid manual override becomes the calculation input.

Changing Liquid Type or Incident Type after making a manual override intentionally loads the standard defaults for the newly selected combination.

## Calculation formulas

The calculator uses the following formulas.

### Total Foam Solution

```text
Total Foam Solution =
    Spill Area × Application Rate × Design Duration
```

With the current units:

```text
sq ft × GPM/sq ft × minutes = gallons
```

### Foam Concentrate Needed

```text
Foam Concentrate Needed =
    Total Foam Solution × (Foam Concentration Percentage / 100)
```

For example, a 3% selection uses a concentration factor of `0.03`.

### Water Needed

```text
Water Needed =
    Total Foam Solution − Foam Concentrate Needed
```

## Example calculation

Assume:

- Hydrocarbon
- Spill Fire
- Spill area = 1,000 sq ft
- Application rate = 0.10 GPM/sq ft
- Design duration = 15 minutes
- Foam concentration = 3%

Total Foam Solution:

```text
1,000 × 0.10 × 15 = 1,500 gallons
```

Foam Concentrate Needed:

```text
1,500 × 0.03 = 45 gallons
```

Water Needed:

```text
1,500 − 45 = 1,455 gallons
```

The calculator therefore displays:

- **Total Foam Solution:** 1,500.00 gal
- **Foam Concentrate Needed:** 45.00 gal
- **Water Needed:** 1,455.00 gal

## User workflow

Open the Calculator page and select the appropriate Liquid Type and Incident Type. Select the foam concentration and enter the spill area.

For the standard workflow, leave Advanced Options closed and press **Calculate**. The application uses the automatically preloaded Application Rate and Design Duration.

To review or alter those parameters, expand **Advanced Options**, edit either field, and press **Calculate**. The calculation uses the values currently shown in the advanced fields.

## Validation

The v6 interface checks that:

- Spill Size is between 100 and 100,000 sq ft.
- Application Rate is greater than zero.
- Design Duration is greater than zero.
- Foam Concentration is a valid percentage.

If an invalid value is detected, the calculator displays an error rather than producing a result.

## Running locally

The project uses Go's standard library and does not require a separate front-end build system.

From the project directory:

```bash
go run . -server
```

By default the application listens on port `8080`.

Open:

```text
http://localhost:8080
```

The server also honors the `PORT` environment variable, which is suitable for hosted services such as Render.

## Project structure

The current application is deliberately compact:

```text
LETSFoamCalculator/
├── assets/
│   └── LOGO.jpg
├── go.mod
├── main.go
└── README.md
```

`main.go` contains the HTTP server, embedded HTML pages, styling, and browser-side calculator logic.

## Deployment

The project repository is:

`https://github.com/letstrain573/LETSFoamCalculator`

The configured Render service is:

`https://letsfoamcalculator.onrender.com`

When a new `main.go` and `README.md` are committed to the deployment branch, Render can rebuild and deploy the updated Go service according to the service's configured build/start commands.

## Version v6 changes

Version v6 adds the functional calculation workflow that was not present in v5:

- Added a Calculate button.
- Added a collapsed-by-default Advanced Options card.
- Added editable Application Rate and Design Duration inputs.
- Added dynamic defaults for all four Liquid Type × Incident Type combinations.
- Added manual override support.
- Added Total Foam Solution calculation.
- Added Foam Concentrate Needed calculation.
- Added Water Needed calculation.
- Added result formatting and input validation.
- Added version display to the calculator page.
- Expanded the home-page description to reflect the calculator's new capabilities.

## Important use note

This calculator is a training and planning aid. Actual firefighting foam selection, application rates, duration, equipment requirements, and tactics should be verified against applicable codes and standards, the foam manufacturer's instructions, the authority having jurisdiction, and the incident-specific conditions before operational use.
