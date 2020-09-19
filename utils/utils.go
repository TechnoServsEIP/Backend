package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-gomail/gomail"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"



var ports = []string{"25576", "25577", "25578", "25579", "26000", "26001"}
var portsBinded = []string{}

func checkBindedPort(port string) bool {
	_, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't listen on port %q: %s\n", port, err)
		return false
	}

	fmt.Printf("TCP Port %q is available", port)
	return true
}

func ReOrderPorts(allPortsBinded []string) {
	portsBinded = allPortsBinded

	if len(portsBinded) > 0 {
		for i := 0; i < len(ports); i++ {
			_, res := Find(portsBinded, ports[i])

			if res {
				ports[i] = ports[len(ports)-1]
				ports[len(ports)-1] = ""
				ports = ports[:len(ports)-1]
				i--
			}
		}
	}
	fmt.Println("ports binded: ", portsBinded)
	fmt.Println("ports: ", ports)
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}, httpCode int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods",
		"GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, "+
			"X-CSRF-Token, Authorization, Origin, X-Auth-Token")
	w.Header().Set("Access-Control-Expose-Headers",
		"Authorization")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(data)
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func FreeThePort(portToFree string) {
	i, res := Find(portsBinded, portToFree)

	if res {
		ports = append(ports, portToFree)

		portsBinded[i] = portsBinded[len(portsBinded)-1]
		portsBinded[len(portsBinded)-1] = ""
		portsBinded = portsBinded[:len(portsBinded)-1]
	}

	fmt.Println("ports binded: ", portsBinded)
	fmt.Println("ports: ", ports)
}

func GetPort() string {
	if len(ports) > 0 {
		portState := checkBindedPort(ports[0])

		if portState {
			/*
			* Get the first port of the ports slice
			 */
			portsBinded = append(portsBinded, ports[0])
			fmt.Println("ports binded: ", portsBinded)

			/*
			* Delete the first port of the ports slice
			 */
			portToSend := ports[0]
			ports = ports[1:len(ports)]
			fmt.Println("ports available: ", ports)

			return portToSend
		} else {
			fmt.Println("no available ports")
			return "no port available"
		}
	}

	return "no port available"
}

func SendConfirmationEmail(url, to string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "jonathan.frickert@epitech.eu")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Confirm Email")
	m.SetBody("text/html", `<!--
==================== Respmail ====================
Respmail is a response HTML email designed to work
on all major devices and responsive for smartphones
that support media queries.
** NOTE **
This template comes with a lot of standard features
that has been thoroughly tested on major platforms
and devices, it is extremely flexible to use and
can be easily customized by removing any row that
you do not need.
it is gauranteed to work 95% without any major flaws,
any changes or adjustments should thoroughly be
tested and reviewed to match with the general
structure.
** Profile **
Licensed under MIT (https://github.com/charlesmudy/responsive-html-email-template/blob/master/LICENSE)
Designed by Shina Charles Memud
Respmail v1.2 (http://charlesmudy.com/respmail/)
** Quick modification **
We are using width of 500 for the whole content,
you can change it any size you want (e.g. 600).
The fastest and safest way is to use find & replace
Sizes: [
		wrapper   : '500',
		columns   : '210',
		x-columns : [
						left : '90',
						right: '350'
				]
		}
	-->

<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="format-detection" content="telephone=no" /> <!-- disable auto telephone linking in iOS -->
    <title>Respmail is a response HTML email designed to work on all major email platforms and smartphones</title>
    <style type="text/css">
        /* RESET STYLES */
        html {
            background-color: #E1E1E1;
            margin: 0;
            padding: 0;
        }

        body,
        #bodyTable,
        #bodyCell,
        #bodyCell {
            height: 100% !important;
            margin: 0;
            padding: 0;
            width: 100% !important;
            font-family: Helvetica, Arial, "Lucida Grande", sans-serif;
        }

        table {
            border-collapse: collapse;
        }

        table[id=bodyTable] {
            width: 100% !important;
            margin: auto;
            max-width: 500px !important;
            color: #7A7A7A;
            font-weight: normal;
        }

        img,
        a img {
            border: 0;
            outline: none;
            text-decoration: none;
            height: auto;
            line-height: 100%;
        }

        a {
            text-decoration: none !important;
            border-bottom: 1px solid;
        }

        h1,
        h2,
        h3,
        h4,
        h5,
        h6 {
            color: #5F5F5F;
            font-weight: normal;
            font-family: Helvetica;
            font-size: 20px;
            line-height: 125%;
            text-align: Left;
            letter-spacing: normal;
            margin-top: 0;
            margin-right: 0;
            margin-bottom: 10px;
            margin-left: 0;
            padding-top: 0;
            padding-bottom: 0;
            padding-left: 0;
            padding-right: 0;
        }

        /* CLIENT-SPECIFIC STYLES */
        .ReadMsgBody {
            width: 100%;
        }

        .ExternalClass {
            width: 100%;
        }

        /* Force Hotmail/Outlook.com to display emails at full width. */
        .ExternalClass,
        .ExternalClass p,
        .ExternalClass span,
        .ExternalClass font,
        .ExternalClass td,
        .ExternalClass div {
            line-height: 100%;
        }

        /* Force Hotmail/Outlook.com to display line heights normally. */
        table,
        td {
            mso-table-lspace: 0pt;
            mso-table-rspace: 0pt;
        }

        /* Remove spacing between tables in Outlook 2007 and up. */
        #outlook a {
            padding: 0;
        }

        /* Force Outlook 2007 and up to provide a "view in browser" message. */
        img {
            -ms-interpolation-mode: bicubic;
            display: block;
            outline: none;
            text-decoration: none;
        }

        /* Force IE to smoothly render resized images. */
        body,
        table,
        td,
        p,
        a,
        li,
        blockquote {
            -ms-text-size-adjust: 100%;
            -webkit-text-size-adjust: 100%;
            font-weight: normal !important;
        }

        /* Prevent Windows- and Webkit-based mobile platforms from changing declared text sizes. */
        .ExternalClass td[class="ecxflexibleContainerBox"] h3 {
            padding-top: 10px !important;
        }

        /* Force hotmail to push 2-grid sub headers down */

        /* /\/\/\/\/\/\/\/\/ TEMPLATE STYLES /\/\/\/\/\/\/\/\/ */

        /* ========== Page Styles ========== */
        h1 {
            display: block;
            font-size: 26px;
            font-style: normal;
            font-weight: normal;
            line-height: 100%;
        }

        h2 {
            display: block;
            font-size: 20px;
            font-style: normal;
            font-weight: normal;
            line-height: 120%;
        }

        h3 {
            display: block;
            font-size: 17px;
            font-style: normal;
            font-weight: normal;
            line-height: 110%;
        }

        h4 {
            display: block;
            font-size: 18px;
            font-style: italic;
            font-weight: normal;
            line-height: 100%;
        }

        .flexibleImage {
            height: auto;
        }

        .linkRemoveBorder {
            border-bottom: 0 !important;
        }

        table[class=flexibleContainerCellDivider] {
            padding-bottom: 0 !important;
            padding-top: 0 !important;
        }

        body,
        #bodyTable {
            background-color: #E1E1E1;
        }

        #emailHeader {
            background-color: #E1E1E1;
        }

        #emailBody {
            background-color: #FFFFFF;
        }

        #emailFooter {
            background-color: #E1E1E1;
        }

        .nestedContainer {
            background-color: #F8F8F8;
            border: 1px solid #CCCCCC;
        }

        .emailButton {
            background-color: #205478;
            border-collapse: separate;
        }

        .buttonContent {
            color: #FFFFFF;
            font-family: Helvetica;
            font-size: 18px;
            font-weight: bold;
            line-height: 100%;
            padding: 15px;
            text-align: center;
        }

        .buttonContent a {
            color: #FFFFFF;
            display: block;
            text-decoration: none !important;
            border: 0 !important;
        }

        .emailCalendar {
            background-color: #FFFFFF;
            border: 1px solid #CCCCCC;
        }

        .emailCalendarMonth {
            background-color: #205478;
            color: #FFFFFF;
            font-family: Helvetica, Arial, sans-serif;
            font-size: 16px;
            font-weight: bold;
            padding-top: 10px;
            padding-bottom: 10px;
            text-align: center;
        }

        .emailCalendarDay {
            color: #205478;
            font-family: Helvetica, Arial, sans-serif;
            font-size: 60px;
            font-weight: bold;
            line-height: 100%;
            padding-top: 20px;
            padding-bottom: 20px;
            text-align: center;
        }

        .imageContentText {
            margin-top: 10px;
            line-height: 0;
        }

        .imageContentText a {
            line-height: 0;
        }

        #invisibleIntroduction {
            display: none !important;
        }

        /* Removing the introduction text from the view */

        /*FRAMEWORK HACKS & OVERRIDES */
        span[class=ios-color-hack] a {
            color: #275100 !important;
            text-decoration: none !important;
        }

        /* Remove all link colors in IOS (below are duplicates based on the color preference) */
        span[class=ios-color-hack2] a {
            color: #205478 !important;
            text-decoration: none !important;
        }

        span[class=ios-color-hack3] a {
            color: #8B8B8B !important;
            text-decoration: none !important;
        }

        /* A nice and clean way to target phone numbers you want clickable and avoid a mobile phone from linking other numbers that look like, but are not phone numbers.  Use these two blocks of code to "unstyle" any numbers that may be linked.  The second block gives you a class to apply with a span tag to the numbers you would like linked and styled.
			Inspired by Campaign Monitor's article on using phone numbers in email: http://www.campaignmonitor.com/blog/post/3571/using-phone-numbers-in-html-email/.
			*/
        .a[href^="tel"],
        a[href^="sms"] {
            text-decoration: none !important;
            color: #606060 !important;
            pointer-events: none !important;
            cursor: default !important;
        }

        .mobile_link a[href^="tel"],
        .mobile_link a[href^="sms"] {
            text-decoration: none !important;
            color: #606060 !important;
            pointer-events: auto !important;
            cursor: default !important;
        }


        /* MOBILE STYLES */
        @media only screen and (max-width: 480px) {

            /*////// CLIENT-SPECIFIC STYLES //////*/
            body {
                width: 100% !important;
                min-width: 100% !important;
            }

            /* Force iOS Mail to render the email at full width. */

            /* FRAMEWORK STYLES */
            /*
				CSS selectors are written in attribute
				selector format to prevent Yahoo Mail
				from rendering media query styles on
				desktop.
				*/
            /*td[class="textContent"], td[class="flexibleContainerCell"] { width: 100%; padding-left: 10px !important; padding-right: 10px !important; }*/
            table[id="emailHeader"],
            table[id="emailBody"],
            table[id="emailFooter"],
            table[class="flexibleContainer"],
            td[class="flexibleContainerCell"] {
                width: 100% !important;
            }

            td[class="flexibleContainerBox"],
            td[class="flexibleContainerBox"] table {
                display: block;
                width: 100%;
                text-align: left;
            }

            /*
				The following style rule makes any
				image classed with 'flexibleImage'
				fluid when the query activates.
				Make sure you add an inline max-width
				to those images to prevent them
				from blowing out.
				*/
            td[class="imageContent"] img {
                height: auto !important;
                width: 100% !important;
                max-width: 100% !important;
            }

            img[class="flexibleImage"] {
                height: auto !important;
                width: 100% !important;
                max-width: 100% !important;
            }

            img[class="flexibleImageSmall"] {
                height: auto !important;
                width: auto !important;
            }


            /*
				Create top space for every second element in a block
				*/
            table[class="flexibleContainerBoxNext"] {
                padding-top: 10px !important;
            }

            /*
				Make buttons in the email span the
				full width of their container, allowing
				for left- or right-handed ease of use.
				*/
            table[class="emailButton"] {
                width: 100% !important;
            }

            td[class="buttonContent"] {
                padding: 0 !important;
            }

            td[class="buttonContent"] a {
                padding: 15px !important;
            }

        }

        /*  CONDITIONS FOR ANDROID DEVICES ONLY
			*   http://developer.android.com/guide/webapps/targeting.html
			*   http://pugetworks.com/2011/04/css-media-queries-for-targeting-different-mobile-devices/ ;
			=====================================================*/

        @media only screen and (-webkit-device-pixel-ratio:.75) {
            /* Put CSS for low density (ldpi) Android layouts in here */
        }

        @media only screen and (-webkit-device-pixel-ratio:1) {
            /* Put CSS for medium density (mdpi) Android layouts in here */
        }

        @media only screen and (-webkit-device-pixel-ratio:1.5) {
            /* Put CSS for high density (hdpi) Android layouts in here */
        }

        /* end Android targeting */

        /* CONDITIONS FOR IOS DEVICES ONLY
			=====================================================*/
        @media only screen and (min-device-width : 320px) and (max-device-width:568px) {}

        /* end IOS targeting */
    </style>
    <!--
			Outlook Conditional CSS
			These two style blocks target Outlook 2007 & 2010 specifically, forcing
			columns into a single vertical stack as on mobile clients. This is
			primarily done to avoid the 'page break bug' and is optional.
			More information here:
			http://templates.mailchimp.com/development/css/outlook-conditional-css
		-->
    <!--[if mso 12]>
			<style type="text/css">
				.flexibleContainer{display:block !important; width:100% !important;}
			</style>
		<![endif]-->
    <!--[if mso 14]>
			<style type="text/css">
				.flexibleContainer{display:block !important; width:100% !important;}
			</style>
		<![endif]-->
</head>

<body bgcolor="#E1E1E1" leftmargin="0" marginwidth="0" topmargin="0" marginheight="0" offset="0">

    <!-- CENTER THE EMAIL // -->
    <!--
		1.  The center tag should normally put all the
			content in the middle of the email page.
			I added "table-layout: fixed;" style to force
			yahoomail which by default put the content left.
		2.  For hotmail and yahoomail, the contents of
			the email starts from this center, so we try to
			apply necessary styling e.g. background-color.
		-->
    <center style="background-color:#E1E1E1;">
        <table border="0" cellpadding="0" cellspacing="0" height="100%" width="100%" id="bodyTable"
            style="table-layout: fixed;max-width:100% !important;width: 100% !important;min-width: 100% !important;">
            <tr>
                <td align="center" valign="top" id="bodyCell">

                    <!-- EMAIL HEADER // -->
                    <!--
							The table "emailBody" is the email's container.
							Its width can be set to 100% for a color band
							that spans the width of the page.
						-->
                    <table bgcolor="#E1E1E1" border="0" cellpadding="0" cellspacing="0" width="500" id="emailHeader">

                        <!-- HEADER ROW // -->
                        <tr>
                            <td align="center" valign="top">
                                <!-- CENTERING TABLE // -->
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td align="center" valign="top">
                                            <!-- FLEXIBLE CONTAINER // -->
                                            <table border="0" cellpadding="10" cellspacing="0" width="500"
                                                class="flexibleContainer">
                                                <tr>
                                                    <td valign="top" width="500" class="flexibleContainerCell">

                                                        <!-- CONTENT TABLE // -->
                                                        <table align="left" border="0" cellpadding="0" cellspacing="0"
                                                            width="100%">
                                                            <tr>
                                                                <!--
																		The "invisibleIntroduction" is the text used for short preview
																		of the email before the user opens it (50 characters max). Sometimes,
																		you do not want to show this message depending on your design but this
																		text is highly recommended.
																		You do not have to worry if it is hidden, the next <td> will automatically
																		center and apply to the width 100% and also shrink to 50% if the first <td>
																		is visible.
																	-->
                                                                <td align="left" valign="middle"
                                                                    id="invisibleIntroduction"
                                                                    class="flexibleContainerBox"
                                                                    style="display:none !important; mso-hide:all;">
                                                                    <table border="0" cellpadding="0" cellspacing="0"
                                                                        width="100%" style="max-width:100%;">
                                                                        <tr>
                                                                            <td align="left" class="textContent">
                                                                                <div
                                                                                    style="font-family:Helvetica,Arial,sans-serif;font-size:13px;color:#828282;text-align:center;line-height:120%;">
                                                                                    The introduction of your message
                                                                                    preview goes here. Try to make it
                                                                                    short.
                                                                                </div>
                                                                            </td>
                                                                        </tr>
                                                                    </table>
                                                                </td>
                                                                <td align="right" valign="middle"
                                                                    class="flexibleContainerBox">
                                                                    <table border="0" cellpadding="0" cellspacing="0"
                                                                        width="100%" style="max-width:100%;">
                                                                        <tr>
                                                                            <td align="left" class="textContent">
                                                                                <!-- CONTENT // -->
                                                                                <div
                                                                                    style="font-family:Helvetica,Arial,sans-serif;font-size:13px;color:#828282;text-align:center;line-height:120%;">
                                                                                    If you can't see this message, <a
                                                                                        href="#" target="_blank"
                                                                                        style="text-decoration:none;border-bottom:1px solid #828282;color:#828282;"><span
                                                                                            style="color:#828282;">view&nbsp;it&nbsp;in&nbsp;your&nbsp;browser</span></a>.
                                                                                </div>
                                                                            </td>
                                                                        </tr>
                                                                    </table>
                                                                </td>
                                                            </tr>
                                                        </table>
                                                    </td>
                                                </tr>
                                            </table>
                                            <!-- // FLEXIBLE CONTAINER -->
                                        </td>
                                    </tr>
                                </table>
                                <!-- // CENTERING TABLE -->
                            </td>
                        </tr>
                        <!-- // END -->

                    </table>
                    <!-- // END -->

                    <!-- EMAIL BODY // -->
                    <!--
							The table "emailBody" is the email's container.
							Its width can be set to 100% for a color band
							that spans the width of the page.
						-->
                    <table bgcolor="#FFFFFF" border="0" cellpadding="0" cellspacing="0" width="500" id="emailBody">

                        <!-- MODULE ROW // -->
                        <!--
								To move or duplicate any of the design patterns
								in this email, simply move or copy the entire
								MODULE ROW section for each content block.
							-->
                        <tr>
                            <td align="center" valign="top">
                                <!-- CENTERING TABLE // -->
                                <!--
										The centering table keeps the content
										tables centered in the emailBody table,
										in case its width is set to 100%.
									-->
                                <table border="0" cellpadding="0" cellspacing="0" width="100%" style="color:#FFFFFF;"
                                    bgcolor="#191E4E">
                                    <tr>
                                        <td align="center" valign="top">
                                            <!-- FLEXIBLE CONTAINER // -->
                                            <!--
													The flexible container has a set width
													that gets overridden by the media query.
													Most content tables within can then be
													given 100% widths.
												-->
                                            <table border="0" cellpadding="0" cellspacing="0" width="500"
                                                class="flexibleContainer">
                                                <tr>
                                                    <td align="center" valign="top" width="500"
                                                        class="flexibleContainerCell">

                                                        <!-- CONTENT TABLE // -->
                                                        <!--
															The content table is the first element
																that's entirely separate from the structural
																framework of the email. 
															-->
                                                        <table border="0" cellpadding="30" cellspacing="0" width="100%">
                                                            <tr>
                                                                <td align="center" valign="top" class="textContent">
                                                                    <h1
                                                                        style="color:#2BCE8B;line-height:100%;font-family:Helvetica,Arial,sans-serif;font-size:35px;font-weight:normal;margin-bottom:5px;text-align:center;">
                                                                        Welcome to Technoservs</h1>
                                                                    <h2
                                                                        style="color: #2BCE8B;text-align:center;font-weight:normal;font-family:Helvetica,Arial,sans-serif;font-size:23px;margin-bottom:10px;color:#205478;line-height:135%;">
                                                                        email confirmation</h2>
                                                                    <div
                                                                        style="text-align:center;font-family:Helvetica,Arial,sans-serif;font-size:15px;margin-bottom:0;color:#FFFFFF;line-height:135%;">
                                                                        This is the confirmation email for your
                                                                        Technoservs acccount if you were not the
                                                                        instigator of this action just ignore it.</div>
                                                                </td>
                                                            </tr>
                                                        </table>
                                                        <!-- // CONTENT TABLE -->

                                                    </td>
                                                </tr>
                                            </table>
                                            <!-- // FLEXIBLE CONTAINER -->
                                        </td>
                                    </tr>
                                </table>
                                <!-- // CENTERING TABLE -->
                            </td>
                        </tr>
                        <!-- // MODULE ROW -->


                        <!-- MODULE ROW // -->
                        <!--  The "mc:hideable" is a feature for MailChimp which allows
								you to disable certain row. It works perfectly for our row structure.
								http://kb.mailchimp.com/article/template-language-creating-editable-content-areas/
							-->
                        <tr mc:hideable>
                            <td align="center" valign="top">
                                <!-- CENTERING TABLE // -->
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td align="center" valign="top">
                                            <!-- FLEXIBLE CONTAINER // -->
                                            <table border="0" cellpadding="30" cellspacing="0" width="500"
                                                class="flexibleContainer">
                                                <tr>
                                                    <td valign="top" width="500" class="flexibleContainerCell">

                                                        <!-- CONTENT TABLE // -->
                                                        <table align="left" border="0" cellpadding="0" cellspacing="0"
                                                            width="100%">
                                                            <tr>
                                                                <td align="center" valign="top"
                                                                    class="flexibleContainerBox">
                                                                    <table border="0" cellpadding="0" cellspacing="0"
                                                                        width="210" style="max-width: 100%;">
                                                                        <tr style="display: flex; margin:auto;">
                                                                            <svg xmlns="http://www.w3.org/2000/svg"
                                                                                xmlns:xlink="http://www.w3.org/1999/xlink"
                                                                                width="442pt" height="294pt"
                                                                                viewBox="0 0 442 294" version="1.1">
                                                                                <g id="surface1">
                                                                                    <path
                                                                                        style="fill:none;stroke-width:30;stroke-linecap:butt;stroke-linejoin:miter;stroke:rgb(9.803922%,11.764706%,30.588235%);stroke-opacity:1;stroke-miterlimit:10;"
                                                                                        d="M 65.189024 82.217736 L 161.677969 82.217736 "
                                                                                        transform="matrix(0.749114,0,0,0.749962,0.138624,0)" />
                                                                                    <path
                                                                                        style="fill:none;stroke-width:30;stroke-linecap:butt;stroke-linejoin:miter;stroke:rgb(9.803922%,11.764706%,30.588235%);stroke-opacity:1;stroke-miterlimit:10;"
                                                                                        d="M 65.981627 147.111672 L 162.470571 147.111672 "
                                                                                        transform="matrix(0.749114,0,0,0.749962,0.138624,0)" />
                                                                                    <path
                                                                                        style=" stroke:none;fill-rule:evenodd;fill:rgb(100%,100%,100%);fill-opacity:1;"
                                                                                        d="M 77.15625 182.796875 L 87.027344 188.503906 L 77.15625 194.21875 Z M 77.15625 182.796875 " />
                                                                                    <path
                                                                                        style=" stroke:none;fill-rule:nonzero;fill:rgb(16.862745%,80.784314%,54.509804%);fill-opacity:1;"
                                                                                        d="M 62.980469 158.21875 L 62.980469 218.792969 L 115.375 188.511719 Z M 62.980469 158.21875 " />
                                                                                    <path
                                                                                        style="fill:none;stroke-width:32;stroke-linecap:butt;stroke-linejoin:miter;stroke:rgb(16.862745%,80.784314%,54.509804%);stroke-opacity:1;stroke-miterlimit:10;"
                                                                                        d="M 570.577506 335.220227 L 571.020738 375.998349 L 216.732566 375.65979 L 24.797576 375.477489 L 24.797576 355.21083 "
                                                                                        transform="matrix(0.749114,0,0,0.749962,0.138624,0)" />
                                                                                    <path
                                                                                        style="fill:none;stroke-width:32;stroke-linecap:butt;stroke-linejoin:miter;stroke:rgb(9.803922%,11.764706%,30.588235%);stroke-opacity:1;stroke-miterlimit:10;"
                                                                                        d="M 561.838019 329.511603 L 562.270822 370.320976 L 207.951363 369.982417 L 16.000729 369.800116 L 16.000729 16.000816 L 562.270822 16.000816 L 562.270822 59.419698 "
                                                                                        transform="matrix(0.749114,0,0,0.749962,0.138624,0)" />
                                                                                    <path
                                                                                        style="fill:none;stroke-width:32;stroke-linecap:butt;stroke-linejoin:miter;stroke:rgb(9.803922%,11.764706%,30.588235%);stroke-opacity:1;stroke-miterlimit:10;"
                                                                                        d="M 208.139085 16.000816 L 207.951363 369.982417 "
                                                                                        transform="matrix(0.749114,0,0,0.749962,0.138624,0)" />
                                                                                    <path
                                                                                        style=" stroke:none;fill-rule:nonzero;fill:rgb(16.862745%,80.784314%,54.509804%);fill-opacity:1;"
                                                                                        d="M 313.753906 77.75 L 313.753906 108.796875 L 274.929688 108.796875 L 274.929688 234.730469 L 237.707031 234.730469 L 237.707031 108.796875 L 198.644531 108.796875 L 198.644531 77.75 Z M 313.753906 77.75 " />
                                                                                    <path
                                                                                        style=" stroke:none;fill-rule:nonzero;fill:rgb(16.862745%,80.784314%,54.509804%);fill-opacity:1;"
                                                                                        d="M 437.613281 112.246094 L 406.367188 122.746094 C 402.167969 110.328125 393.265625 104.113281 379.664062 104.109375 C 365.144531 104.109375 357.886719 108.828125 357.886719 118.261719 C 357.828125 121.878906 359.3125 125.347656 361.960938 127.808594 C 364.671875 130.414062 370.808594 132.675781 380.375 134.59375 C 396.355469 137.824219 408.054688 141.085938 415.476562 144.382812 C 422.898438 147.675781 429.140625 152.988281 434.207031 160.3125 C 439.261719 167.492188 441.941406 176.074219 441.867188 184.859375 C 441.867188 199.042969 436.414062 211.308594 425.507812 221.660156 C 414.601562 232.007812 398.4375 237.179688 377.019531 237.175781 C 360.953125 237.175781 347.273438 233.515625 335.980469 226.195312 C 324.691406 218.875 317.257812 207.996094 313.6875 193.558594 L 347.699219 185.855469 C 351.519531 200.5 362.140625 207.820312 379.566406 207.820312 C 387.980469 207.820312 394.246094 206.191406 398.367188 202.933594 C 402.488281 199.671875 404.554688 195.78125 404.5625 191.261719 C 404.710938 186.949219 402.574219 182.878906 398.9375 180.5625 C 395.191406 178.03125 387.964844 175.652344 377.257812 173.429688 C 357.28125 169.289062 343.046875 163.632812 334.558594 156.464844 C 326.066406 149.292969 321.824219 138.808594 321.824219 125.011719 C 321.824219 111.050781 326.890625 99.28125 337.03125 89.703125 C 347.167969 80.125 360.917969 75.34375 378.285156 75.363281 C 409.039062 75.363281 428.8125 87.65625 437.613281 112.246094 Z M 437.613281 112.246094 " />
                                                                                    <path
                                                                                        style=" stroke:none;fill-rule:nonzero;fill:rgb(9.803922%,11.764706%,30.588235%);fill-opacity:1;"
                                                                                        d="M 306.667969 73.488281 L 306.667969 104.542969 L 267.84375 104.542969 L 267.84375 230.476562 L 230.617188 230.476562 L 230.617188 104.542969 L 191.582031 104.542969 L 191.582031 73.488281 Z M 306.667969 73.488281 " />
                                                                                    <path
                                                                                        style=" stroke:none;fill-rule:nonzero;fill:rgb(9.803922%,11.764706%,30.588235%);fill-opacity:1;"
                                                                                        d="M 430.550781 107.996094 L 399.304688 118.457031 C 395.105469 106.035156 386.203125 99.828125 372.597656 99.828125 C 358.082031 99.828125 350.820312 104.542969 350.820312 113.972656 C 350.761719 117.585938 352.242188 121.058594 354.890625 123.519531 C 357.601562 126.128906 363.738281 128.394531 373.300781 130.3125 C 389.285156 133.539062 400.984375 136.800781 408.40625 140.101562 C 415.828125 143.398438 422.070312 148.707031 427.132812 156.023438 C 432.179688 163.207031 434.851562 171.792969 434.765625 180.574219 C 434.765625 194.761719 429.3125 207.023438 418.40625 217.367188 C 407.5 227.714844 391.335938 232.886719 369.917969 232.894531 C 353.84375 232.894531 340.167969 229.230469 328.878906 221.90625 C 317.59375 214.582031 310.152344 203.714844 306.554688 189.304688 L 340.566406 181.597656 C 344.386719 196.238281 355.007812 203.5625 372.433594 203.5625 C 380.84375 203.5625 387.109375 201.933594 391.234375 198.679688 C 395.363281 195.425781 397.425781 191.53125 397.429688 187.003906 C 397.582031 182.691406 395.441406 178.625 391.804688 176.308594 C 388.058594 173.785156 380.832031 171.40625 370.125 169.175781 C 350.148438 165.027344 335.917969 159.371094 327.425781 152.203125 C 318.9375 145.039062 314.691406 134.554688 314.691406 120.75 C 314.691406 106.800781 319.761719 95.035156 329.898438 85.449219 C 340.035156 75.867188 353.789062 71.074219 371.160156 71.074219 C 401.9375 71.074219 421.734375 83.378906 430.550781 107.996094 Z M 430.550781 107.996094 " />
                                                                                </g>
                                                                            </svg>

                                                                        </tr>
                                                                    </table>
                                                                </td>
                                                                <td align="right" valign="middle"
                                                                    class="flexibleContainerBox">
                                                                    <table class="flexibleContainerBoxNext" border="0"
                                                                        cellpadding="0" cellspacing="0" width="210"
                                                                        style="max-width: 100%;">
                                                                        <tr>
                                                                            <td align="left" class="textContent">
                                                                                <div
                                                                                    style="text-align:left;font-family:Helvetica,Arial,sans-serif;font-size:15px;margin-bottom:0;color:#5F5F5F;line-height:135%;">
                                                                                </div>
                                                                            </td>
                                                                        </tr>
                                                                    </table>
                                                                </td>
                                                            </tr>
                                                        </table>
                                                        <!-- // CONTENT TABLE -->

                                                    </td>
                                                </tr>
                                            </table>
                                            <!-- // FLEXIBLE CONTAINER -->
                                        </td>
                                    </tr>
                                </table>
                                <!-- // CENTERING TABLE -->
                            </td>
                        </tr>
                        <!-- // MODULE ROW -->


                        <!-- MODULE ROW // -->
                        <tr>
                            <td align="center" valign="top">
                                <!-- CENTERING TABLE // -->
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr style="padding-top:0;">
                                        <td align="center" valign="top">
                                            <!-- FLEXIBLE CONTAINER // -->
                                            <table border="0" cellpadding="30" cellspacing="0" width="500"
                                                class="flexibleContainer">
                                                <tr>
                                                    <td style="padding-top:0;" align="center" valign="top" width="500"
                                                        class="flexibleContainerCell">

                                                        <!-- CONTENT TABLE // -->
                                                        <table border="0" cellpadding="0" cellspacing="0" width="50%"
                                                            class="emailButton" style="background-color: #2BCE8B;">
                                                            <tr>
                                                                <td align="center" valign="middle" class="buttonContent"
                                                                    style="padding-top:15px;padding-bottom:15px;padding-right:15px;padding-left:15px;">
                                                                    <a style="color:#FFFFFF;text-decoration:none;font-family:Helvetica,Arial,sans-serif;font-size:20px;line-height:135%;"
                                                                        href="`+url+`" target="_blank">Confirmation</a>
                                                                </td>
                                                            </tr>
                                                        </table>
                                                        <!-- // CONTENT TABLE -->`)

	d := gomail.NewDialer("in-v3.mailjet.com", 587, "56a430fb737fca0b6c5d33a449c6206e", "097628559f09a9bab73a6fab8b2d357d")

	return d.DialAndSend(m)
}

func SendInvitationEmail(sender, url, to string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "jonathan.frickert@epitech.eu")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Invitation on a minecraft server")
	m.SetBody("text/html", ""+sender+" send you a minecraft server url: "+url+"")

	d := gomail.NewDialer("in-v3.mailjet.com", 587, "56a430fb737fca0b6c5d33a449c6206e", "097628559f09a9bab73a6fab8b2d357d")

	return d.DialAndSend(m)
}

//func DecryptToken(token string) (jwt.StandardClaims, bool, error)  {
//	decryptToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, errors.New("error")
//		}
//		return []byte(os.Getenv("token_password")), nil
//	})
//	claims, valid := decryptToken.Claims.(jwt.StandardClaims)
//	if !decryptToken.Valid {
//		valid = false
//	}
//	return claims, valid, err
//}

func DecryptToken(tokenString string) (jwt.StandardClaims, bool, error) {
	decryptToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})
	if err != nil { //Malformed token, returns with http code 403 as usual
		fmt.Println("Malformed authentication token", err)
		return jwt.StandardClaims{}, false, errors.New("error")
	}
	claims, valid := decryptToken.Claims.(jwt.StandardClaims)
	if !decryptToken.Valid {
		valid = false
	}
	return claims, valid, nil

	//tokenParsed, err := jwt.ParseWithClaims(tokenString, jwt.StandardClaims{},
	//func(token *jwt.Token) (interface{}, error) {
	//	return []byte(os.Getenv("token_password")), nil
	//})
	//if err != nil { //Malformed token, returns with http code 403 as usual
	//	fmt.Println("Malformed authentication token")
	//	return tokenParsed.Claims.(jwt.StandardClaims), false, errors.New("error")
	//}
	//return tokenParsed.Claims.(jwt.StandardClaims), tokenParsed.Valid, nil
}

func CreateTmpFolder(folderName string) error {
	cmd := exec.Command("mkdir", folderName)

	_, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	return nil
}

func DeleteTmpFolder(folderName string) error {
	cmd := exec.Command("rm", "-rf", folderName)

	_, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	return nil
}

//  Example:
//  local to container
//  path => {localPath}/{file}, destination => {containerID}:/{path}/{file}
//  container to local
//  path => {containerID}:/{path}/{file}, destination => {localPath}/{file}
func DockerCopy(path string, destination string) error {
	cmd := exec.Command("docker", "cp", path, destination)

	_, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	return nil
}
