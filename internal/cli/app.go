package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"example.com/xiangzhenfarm/internal/domain"
	"example.com/xiangzhenfarm/internal/importer"
	"example.com/xiangzhenfarm/internal/report"
	"example.com/xiangzhenfarm/internal/service"
)

type App struct {
	Registry *service.Registry
	In       io.Reader
	Out      io.Writer
}

func New(registry *service.Registry, input io.Reader, output io.Writer) *App {
	return &App{Registry: registry, In: input, Out: output}
}

func (a *App) Run(args []string) error {
	if len(args) == 0 {
		a.printHelp()
		return nil
	}
	switch args[0] {
	case "add":
		return a.add(args[1:])
	case "list":
		return a.list(args[1:])
	case "search":
		return a.search(args[1:])
	case "update":
		return a.update(args[1:])
	case "delete":
		return a.delete(args[1:])
	case "import":
		return a.importCSV(args[1:])
	case "report":
		return a.summary()
	case "archive":
		return a.archive(args[1:])
	case "restore":
		return a.restore(args[1:])
	case "history":
		return a.history(args[1:])
	case "dashboard":
		return a.dashboard()
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func (a *App) add(args []string) error {
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	head := flags.String("head", "", "household head")
	village := flags.String("village", "", "village group")
	area := flags.Float64("area", 0, "cultivated area")
	crop := flags.String("crop", "", "main crop")
	phone := flags.String("phone", "", "phone")
	visit := flags.String("visit", "", "last visit")
	if err := flags.Parse(args); err != nil {
		return err
	}
	record, err := a.Registry.CreateFarmer(domain.FarmerRecord{HouseholdHead: *head, VillageGroup: *village, CultivatedArea: *area, MainCrop: *crop, Phone: *phone, LastVisit: *visit})
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.Out, report.RecordLine(record))
	return err
}

func (a *App) list(args []string) error {
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	field := flags.String("sort", "village_group", "sort field")
	reverse := flags.Bool("reverse", false, "reverse order")
	if err := flags.Parse(args); err != nil {
		return err
	}
	records, err := a.Registry.SearchFarmers(domain.SearchFilter{}, domain.SortSpec{Field: *field, Reverse: *reverse})
	if err != nil {
		return err
	}
	return printRecords(a.Out, records)
}

func (a *App) search(args []string) error {
	flags := flag.NewFlagSet("search", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	crop := flags.String("crop", "", "crop")
	village := flags.String("village", "", "village group")
	head := flags.String("head", "", "household head")
	if err := flags.Parse(args); err != nil {
		return err
	}
	records, err := a.Registry.SearchFarmers(domain.SearchFilter{MainCrop: *crop, VillageGroup: *village, HouseholdHead: *head}, domain.SortSpec{Field: "village_group"})
	if err != nil {
		return err
	}
	return printRecords(a.Out, records)
}

func (a *App) update(args []string) error {
	if len(args) < 1 {
		return errors.New("update requires farmer id")
	}
	flags := flag.NewFlagSet("update", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	phone := flags.String("phone", "", "phone")
	crop := flags.String("crop", "", "main crop")
	area := flags.Float64("area", 0, "cultivated area")
	visit := flags.String("visit", "", "last visit")
	notes := flags.String("notes", "", "notes")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	record, err := a.Registry.UpdateFarmer(args[0], domain.FarmerRecord{Phone: *phone, MainCrop: *crop, CultivatedArea: *area, LastVisit: *visit, Notes: *notes})
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.Out, report.RecordLine(record))
	return err
}

func (a *App) delete(args []string) error {
	if len(args) != 1 {
		return errors.New("delete requires farmer id")
	}
	if err := a.Registry.DeleteFarmer(args[0]); err != nil {
		return err
	}
	_, err := fmt.Fprintln(a.Out, "deleted="+args[0])
	return err
}

func (a *App) importCSV(args []string) error {
	if len(args) != 1 {
		return errors.New("import requires a csv path")
	}
	file, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer file.Close()
	result, err := importer.ProcessFiles([]importer.BatchFile{{Name: args[0], Input: file}}, a.Registry)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(a.Out, "files=%d imported=%d rejected=%d\n", result.Files, result.Imported, result.Rejected)
	return err
}

func (a *App) summary() error {
	summary, err := report.Build(a.Registry)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.Out, report.Render(summary))
	return err
}

func (a *App) archive(args []string) error {
	if len(args) < 4 {
		return errors.New("archive requires id, operator, date, and reason")
	}
	entry, err := a.Registry.ArchiveFarmer(args[0], args[1], args[2], strings.Join(args[3:], " "))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(a.Out, "archived=%s\n", entry.FarmerID)
	return err
}

func (a *App) restore(args []string) error {
	if len(args) != 1 {
		return errors.New("restore requires farmer id")
	}
	record, err := a.Registry.RestoreFarmer(args[0])
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.Out, report.RecordLine(record))
	return err
}

func (a *App) history(args []string) error {
	if len(args) != 1 {
		return errors.New("history requires farmer id")
	}
	events, err := a.Registry.FarmerTimeline(args[0])
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.Out, report.RenderTimeline(events))
	return err
}

func (a *App) dashboard() error {
	dashboard, err := a.Registry.BuildDashboard()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(a.Out, dashboard.Render())
	return err
}

func (a *App) printHelp() {
	_, _ = fmt.Fprintln(a.Out, "farmregistry add|list|search|update|delete|import|report|archive|restore|history|dashboard")
}

func printRecords(output io.Writer, records []domain.FarmerRecord) error {
	for _, record := range records {
		if _, err := fmt.Fprintln(output, report.RecordLine(record)); err != nil {
			return err
		}
	}
	return nil
}

func ParseArea(value string) (float64, error) {
	area, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil || area <= 0 {
		return 0, errors.New("area must be a positive number")
	}
	return area, nil
}
