
function addColorText(layer, index, color) {
    var animator = layer.property('Text').property('Animators').addProperty('ADBE Text Animator');
    var colorProp = animator.property('Properties').addProperty('ADBE Text Fill Color');
    colorProp.setValue(color);
    var selector = animator.property('Selectors').addProperty('ADBE Text Selector');
    var advanced = selector.property('Advanced');
    advanced.property('Based On').setValue(2); // Characters
    advanced.property('Units').setValue(2); // Index
    var start = selector.property('ADBE Text Index Start').setValue(index)
    var end = selector.property('ADBE Text Index End').setValue(index+1)
}


var font_size = 30;
var tracking = 700;

var positionX = 960;
var positionY = 540;

var width = 1800;
var height = 900;

var letterWidth = 22.2
var letterHeight = 36.0
var wrapLength = 22

var originX = positionX - width/2
var originY = positionY - height/2

var text = "DIFNOTAHOPELESSONENORCOULDANYCREWENDUREFORSOLONGAPERIODSUCHANUNINTERMITTEDINTENSESTRAININGATTHEOARATHINGBARELYTOLERABLEONLYINSOMEONEBRIEFVICISSITUDETHESHIPITSELFTHENASITSOMETIMESHAPPENSOFFEREDTHEMOSTPROMISINGINTERMEDIATEMEANSOFOVERTAKINGTHECHASEACCORDINGLYTHEBOATSNOWMADEFORHERANDWERESOONSWAYEDUPTOTHEIRCRANESTHETWOPARTSOFTHEWRECKEDBOATHAVINGBEENPREVIOUSLYSECUREDBYHERANDTHENHOISTINGEVERYTHINGTOHERSIDEANDSTACKINGHERCANVASHIGHUPANDSIDEWAYSOUTSTRETCHINGITWITHSTUNSAILSLIKETHEDOUBLEJOINTEDWINGSOFANALBATROSSTHEPEQUODBOREDOWNINTHELEEWARDWAKEOFMOBYDICKATTHEWELLKNOWNMETHODICINTERVALSTHEWHALESGLITTERINGSPOUTWASREGULARLYANNOUNCEDFROMTHEMANNEDMASTHEADSANDWHENHEWOULDBEREPORTEDASJUSTGONEDOWNAHABWOULDTAKETHETIMEANDTHENPACINGTHEDECKBINNACLEWATCHINHANDSOSOONASTHELASTSECONDOFTHEALLOTTEDHOUREXPIREDHISVOICEWASHEARDWHOSEISTHEDOUBLOONNOWDYESEEHIMANDIFTHEREPLYWASNOSIRSTRAIGHTWAYHECOMMANDEDTHEMTOLIFTHIMTOHISPERCHINTHISWAYTHEDAYWOREONAHABNOWALOFTANDMOTIONLESSANONUNRESTINGLYPACINGTHEPLANKSASHEWASTHUSWALKINGUTTERINGNOSOUNDEXCEPTTOHAILTHEMENALOFTORTOBIDTHEMHOISTASAILSTILLHIGHERORTOSPREADONETOASTILLGREATERBREADTHTHUSTOANDFROPACINGBENEATHHISSLOUCHEDHATATEVERYTURNHEPASSEDHISOWNWRECKEDBOATWHICHHADBEENDROPPEDUPONTHEQUARTERDECKANDLAYTHEREREVERSEDBROKENBOWTOSHATTEREDSTERNATLASTHEPAUSEDBEFOREITANDASINANALREADYOVERCLOUDEDSKYFRESHTROOPSOFCLOUDSWILLSOMETIMESSAILACROSSSOOVERTHEOLDMANSFACETHERENOWSTOLESOMESUCHADDEDGLOOMASTHISSTUBBSAWHIMPAUSEANDPERHAPSINTENDINGNOTVAINLYTHOUGHTOEVINCEHISOWNUNABATEDFORTITUDEANDTHUSKEEPUPAVALIANTPLACEINHISCAPTAINSMINDHEADVANCEDANDEYEINGTHEWRECKEXCLAIMEDTHETHISTLETHEASSREFUSEDITPRICKEDHISMOUTHTOOKEENLYSIRHAHAWHATSOULLESSTHINGISTHISTHATLAUGHSBEFOREAWRECKMANMANDIDINOTKNOWTHEEBRAVEASFEARLESSFIREANDASMECHANICALICOULDSWEARTHOUWERTAPOLTROONGROANNORLAUGHSHOULDBEHEARDBEFOREAWRECKAYESIRSAIDSTARBUCKDRAWINGNEARTISASOLEMNSIGHTANOMENANDANILLONEOMENOMENTHEDICTIONARYIFTHEGODSTHINKTOSPEAKOUTRIGHTTOMANTHEYWILLHONORABLYSPEAKOUTRIGHTNOTSHAKETHEIRHEADSANDGIVEANOLDWIVESDARKLINGHINTBEGONEYETWOARETHEOPPOSITEPOLESOFONETHINGSTARBUCKISSTUBBREVERSEDANDSTUBBISSTARBUCKANDYETWOAREALLMANKINDANDAHABSTANDSALONEAMONGTHEMILLIONSOFTHEPEOPLEDEARTHNORGODSNORMENHISNEIGHBORSCOLDCOLDISHIVERHOWNOWALOFTTHEREDYESEEHIMSINGOUTFOREVERYSPOUTTHOUGHHESPOUTTENTIMESASECONDTHEDAYWASNEARLYDONEONLYTHEHEMOFHISGOLDENROBEWASRUSTLINGSOONITWASALMOSTDARKBUTTHELOOKOUTMENSTILLREMAINEDUNSETCANTSEETHESPOUTNOWSIRTOODARKCRIEDAVOICEFROMTHEAIRHOWHEADINGWHENLASTSEENASBEFORESIRSTRAIGHTTOLEEWARDGOODHEWILLTRAVELSLOWERNOWTISNIGHTDOWNROYALSANDTOPGALLANTSTUNSAILSMRSTARBUCKWEMUSTNOTRUNOVERHIMBEFOREMORNINGHESMAKINGAPASSAGENOWANDMAYHEAVETOAWHILEHELMTHEREKEEPHERFULLBEFORETHEWINDALOFTCOMEDOWNMRSTUBBSENDAFRESHHANDTOTHEFOREMASTHEADANDSEEITMANNEDTILLMORNINGTHENADVANCINGTOWARDSTHEDOUBLOONINTHEMAINMASTMENTHISGOLDISMINEFORIEARNEDITBUTISHALLLETITABIDEHERETILLTHEWHITEWHALEISDEADANDTHENWHOSOEVEROFYEFIRSTRAISESHIMUPONTHEDAYHESHALLBEKILLEDTHISGOLDISTHATMANSANDIFONTHATDAYISHALLAGAINRAISEHIMTHENTENTIMESITSSUMSHALLBEDIVIDEDAMONGALLOFYEAWAYNOWTHEDECKISTHINESIRANDSOSAYINGHEPLACEDHIMSELFHALFWAYWITHINTHESCUTTLEANDSLOUCHINGHISHATSTOODTHERETILLDAWNEXCEPTWHENATINTERVALSROUSINGHIMSELFTOSEEHOWTHENIGHTWOREONTHISMOTIONISPECULIARTOTHESPERMWHALEITRECEIVESITSDESIGNATIONPITCHPOLINGFROMITSBEINGLIKENEDTOTHATPRELIMINARYUPANDDOWNPOISEOFTHEWHALELANCEINTHEEXERCISECALLEDPITCHPOLINGPREVIOUSLYDESCRIBEDBYTHISMOTIONTHEWHALEMUSTBESTANDMOSTCOMPREHENSIVELYVIEWWHATEVEROBJECTSMAYBEENCIRCLINGHIMATDAYBREAKTHETHREEMASTHEADSWEREPUNCTUALLYMANNEDAFRESHDYESEEHIMCRIEDAHABAFTERALLOWINGALITTLESPACEFORTHELIGHTTOSPREADSEENOTHINGSIRTURNUPALLHANDSANDMAKESAILHETRAVELSFASTERTHANITHOUGHTFORTHETOPGALLANTSAILSAYETHEYSHOULDHAVEBEENKEPTONHERALLNIGHTBUTNOMATTERTISBUTRESTINGFORTHERUSHHEREBEITSAIDTHATTHISPERTINACIOUSPURSUITOFONEPARTICULARWHALECONTINUEDTHROUGHDAYINTONIGHTANDTHROUGHNIGHTINTODAYISATHINGBYNOMEANSUNPRECEDENTEDINTHESOUTHSEAFISHERYFORSUCHISTHEWONDERFULSKILLPRESCIENCEOFEXPERIENCEANDINVINCIBLECONFIDENCEACQUIREDBYSOMEGREATNATURALGENIUSESAMONGTHENANTUCKETCOMMANDERST";

app.beginUndoGroup('bible-code');
var comp = app.project.activeItem;

label = comp.layers.addBoxText([width,height], text);
label.threeDLayer = true

var textProp = label.property("Source Text");
label.position.setValue([positionX, positionY]);
var textDocument = textProp.value;
textDocument.fontSize = font_size;
textDocument.fillColor = [0, 0, 0];
textDocument.tracking = tracking;
textProp.setValue(textDocument);

// "NOERTAP": 889914-891108 skip:199 
addColorText(label, 2256, [0.9, 0.0, 0.0])  //N
addColorText(label, 2455, [0.9, 0.0, 0.0])  //O
addColorText(label, 2654, [0.9, 0.0, 0.0])  //E
addColorText(label, 2853, [0.9, 0.0, 0.0])  //R
addColorText(label, 3052, [0.9, 0.0, 0.0])  //T
addColorText(label, 3251, [0.9, 0.0, 0.0])  //A
addColorText(label, 3450, [0.9, 0.0, 0.0])  //P

// "DATA": 888534-889746 skip:404 
addColorText(label, 876, [0.0, 0.9, 0.5])  //D
addColorText(label, 1280, [0.0, 0.9, 0.5])  //A
addColorText(label, 1684, [0.0, 0.9, 0.5])  //T
addColorText(label, 2088, [0.0, 0.9, 0.5])  //A

// "EMIT": 889907-891098 skip:397 
addColorText(label, 2249, [0.0, 0.0, 0.9])  //E
addColorText(label, 2646, [0.0, 0.0, 0.9])  //M
addColorText(label, 3043, [0.0, 0.0, 0.9])  //I
addColorText(label, 3440, [0.0, 0.0, 0.9])  //T



app.endUndoGroup();