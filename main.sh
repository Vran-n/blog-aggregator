echo "Bootdev wrapper. Commands:" 
echo "\"-r\"        run" 
echo "\"-s\"        submit" 
echo "\"--up\"      up database migration" 
echo "\"--down\"    down database migration" 

read OPTION

case $OPTION in
    -r)
        echo ~~ DATABASE UP MIGRATION
        ./scripts/db_up.sh
        bootdev run
        echo ~~ DATABASE DOWN MIGRATION
        ./scripts/db_down.sh
        ;;
    -s)
        echo ~~ DATABASE UP MIGRATION
        ./scripts/db_up.sh
        bootdev run -s
        echo ~~ DATABASE DOWN MIGRATION
        ./scripts/db_down.sh
        ;;
    --up)
        echo ~~ DATABASE UP MIGRATION
        ./scripts/db_up.sh
        ;;
    --down)
        echo ~~ DATABASE DOWN MIGRATION
        ./scripts/db_down.sh
        ;;
    *)
        echo "what the fuck \"$OPTION\" is this piece of shit"
        ;;
esac