//go:generate   go-bindata-assetfs -prefix ../client/app/  ../client/app/auth/... ../client/app/clients/... ../client/app/modules/... ../client/app/module/... ../client/app/css/... ../client/app/fonts/... ../client/app/libs/... ../client/app/reports/... ../client/app/users/... ../client/app/
package main

type Config struct{}

func main() {
	/*svcConfig := &service.Config{
		Name:        "fakelive",
		DisplayName: "fakelive and opinion server",
		Description: "",
	}*/

	prg := &app{}
	prg.run()
	select {}
	/*s, err := service.New(prg, svcConfig)
	       if err != nil {
		               log.Fatal(err)
		       }
	       if len(os.Args) > 1 {
		               err = service.Control(s, os.Args[1])
		               if err != nil {
			                       log.Fatal(err)
			               }
		               return
		       }

	       logger, err := s.Logger(nil)
	       if err != nil {
		               log.Fatal(err)
		       }
	       err = s.Run()
	       if err != nil {
		               logger.Error(err)
		       }*/

}
